//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/michaelprice232/user-mgmt-service-api/tests/endpoints"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	httpMaxRetries         = 5
	httpTimeBetweenRetries = 2 * time.Second

	FargateTaskTimeBetweenRetries = 20 * time.Second
	FargateTaskTimeout            = 10 * time.Minute
)

func TestUsingAWS(t *testing.T) {
	// Pull the Docker image refs from envars so we can pass the branch builds from the CI system
	// Assumes the DB seeder image will always be based on the same Docker tag as the app image
	// Use Terraform defaults if not present, such as running locally
	var varInput map[string]interface{}
	appImage := os.Getenv("DOCKER_APP_IMAGE")
	if appImage != "" {
		dbSeedImage := fmt.Sprintf("%s-db-seeding", appImage)
		varInput = map[string]interface{}{
			"unique_identifier":    strings.ToLower(random.UniqueId()),
			"fargate_docker_image": appImage,
			"e2e_db_seed_image":    dbSeedImage,
		}
		t.Logf("Using docker app image '%s' and DB seeder image '%s'", appImage, dbSeedImage)
	} else {
		varInput = map[string]interface{}{
			"unique_identifier": strings.ToLower(random.UniqueId()),
		}
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../terraform",
		Vars:         varInput,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get outputs
	baseURL := terraform.Output(t, terraformOptions, "service_endpoint")
	ecsClusterName := terraform.Output(t, terraformOptions, "ecs_cluster_name")
	targetSubnet := terraform.Output(t, terraformOptions, "db_seeding_target_subnet")
	securityGroup := terraform.Output(t, terraformOptions, "fargate_task_security_group_id")
	targetDefinitions := terraform.Output(t, terraformOptions, "db_seeding_task_definition_target")
	serviceName := terraform.Output(t, terraformOptions, "service_name")

	// Create ECS client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
		t.Fatalf("unable to load SDK config: %v", err.Error())
	}
	ecsClient := ecs.NewFromConfig(cfg)

	seedDatabase(t, ecsClient, ecsClusterName, targetSubnet, securityGroup, targetDefinitions)

	validateService(t, ecsClient, ecsClusterName, serviceName)

	endpoints.CheckEndpoints(t, baseURL, httpMaxRetries, httpTimeBetweenRetries)
}

func validateService(t *testing.T, ecsClient *ecs.Client, clusterName, serviceName string) {
	t.Logf("Waiting for service %s to be ready", serviceName)

	ctx, cancelFunc := context.WithTimeout(context.Background(), FargateTaskTimeout)
	defer cancelFunc()
	currentAttempt := 1

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Waiting for service to be ready timed out after %v: %v", FargateTaskTimeout, ctx.Err())
		default:
			ready := true

			describeResult, err := ecsClient.DescribeServices(context.Background(), &ecs.DescribeServicesInput{
				Cluster:  aws.String(clusterName),
				Services: []string{serviceName},
			})
			if err != nil {
				t.Fatalf("Failed to describe services: %v", err.Error())
			}

			if len(describeResult.Services) != 1 {
				t.Logf("Waiting for service %s to be ready. Expected 0 services but got %d", serviceName, len(describeResult.Services))
				ready = false
			}

			service := describeResult.Services[0]

			if *service.Status != "ACTIVE" {
				t.Logf("Waiting for service %s to be ready. Expected service to the ACTIVE but got %s", serviceName, *service.Status)
				ready = false

			}

			if service.DesiredCount != service.RunningCount {
				t.Logf("Waiting for service %s to be ready. Expected the running replicas (%d) to match the desired replicas (%d)", serviceName, service.RunningCount, service.DesiredCount)
				ready = false
			}

			if ready {
				return
			}

			t.Logf("Service not ready yet. Attempt: %d, sleeping for %s", currentAttempt, FargateTaskTimeBetweenRetries)
			time.Sleep(FargateTaskTimeBetweenRetries)
			currentAttempt++
		}
	}
}

// seedDatabase prepares the database for E2E testing by creating a table and some sample data.
func seedDatabase(t *testing.T, ecsClient *ecs.Client, ecsClusterName, targetSubnet, securityGroup, targetDefinitions string) {
	// Deploy Fargate task to update database
	ecsResult, err := ecsClient.RunTask(context.Background(), &ecs.RunTaskInput{
		TaskDefinition: aws.String(targetDefinitions),
		Cluster:        aws.String(ecsClusterName),
		Count:          aws.Int32(1),
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{targetSubnet},
				AssignPublicIp: types.AssignPublicIpDisabled,
				SecurityGroups: []string{securityGroup},
			},
		},
	})
	assert.NoError(t, err, "Expected no errors from running the ECS task")
	assert.Empty(t, ecsResult.Failures, "Expected no error results from running the ECS task")

	// Wait until the DB seeding task has completed (with retries and a timeout)
	taskID := *ecsResult.Tasks[0].TaskArn
	t.Logf("Waiting for DB seeding task %s to complete", taskID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), FargateTaskTimeout)
	defer cancelFunc()
	currentAttempt := 1

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("DB seed task timed out after %v: %v", FargateTaskTimeout, ctx.Err())
		default:
			describeResult, err := ecsClient.DescribeTasks(context.Background(), &ecs.DescribeTasksInput{
				Cluster: aws.String(ecsClusterName),
				Tasks:   []string{taskID},
			})
			if err != nil {
				t.Fatalf("Failed to describe tasks: %v", err.Error())
			}

			// Wait until the task is in a stopped state with a successful exit code
			if describeResult.Tasks[0].LastStatus != nil && *describeResult.Tasks[0].LastStatus == "STOPPED" {
				if *describeResult.Tasks[0].Containers[0].ExitCode == 0 {
					t.Logf("DB seeding task completed successfully")
					break
				} else {
					t.Fatalf("DB seeding task completed with wrong exit code. Expected 0, got %v", *describeResult.Tasks[0].Containers[0].ExitCode)
				}
			}

			t.Logf("Task status currently: %s. Attempt: %d, sleeping for %s", *describeResult.Tasks[0].LastStatus, currentAttempt, FargateTaskTimeBetweenRetries)
			time.Sleep(FargateTaskTimeBetweenRetries)
			currentAttempt++
		}
	}
}

// todo: health check version does not match branch
