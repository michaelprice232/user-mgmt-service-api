# Bare-bones Terraform (for E2E tests)

`Non-production ready` Terraform code so we can run some E2E tests from the pipeline.
For example, no TLS, WAF, custom DNS zones or AuthN/AuthZ is used.

The database (AWS Aurora) also requires seeding with the SQL table and data before the E2E tests can be run.
This would typically be performed by a tool such as [Liquidbase](https://www.liquibase.com/) in production, but I am doing it as part of the 
Terratest E2E tests for brevity.