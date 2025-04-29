terraform {
  backend "s3" {
    bucket         = "smart-retention-terraform-state"
    key            = "global/smart-retention/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    use_lockfile   = true
  }
}
