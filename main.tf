provider "aws" {
  access_key = "test"
  secret_key = "test"
  region     = "us-east-1"
  s3_use_path_style           = false // Only required for non virtual hosted-style endpoint use case. https://registry.terraform.io/providers/hashicorp/aws/latest/docs#s3_force_path_style
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway          = "http://localhost:4566"
    apigatewayv2        = "http://localhost:4566"
    cloudformation      = "http://localhost:4566"
    cloudwatch          = "http://localhost:4566"
    cloudwatchlogs      = "http://localhost:4566"
    cloudwatchevents    = "http://localhost:4566"
    dynamodb            = "http://localhost:4566"
    ec2                 = "http://localhost:4566"
    es                  = "http://localhost:4566"
    elasticache         = "http://localhost:4566"
    eventbridge         = "http://localhost:4566"
    firehose            = "http://localhost:4566"
    iam                 = "http://localhost:4566"
    kinesis             = "http://localhost:4566"
    lambda              = "http://localhost:4566"
    rds                 = "http://localhost:4566"
    redshift            = "http://localhost:4566"
    route53             = "http://localhost:4566"
    s3                  = "http://s3.localhost.localstack.cloud:4566"
    secretsmanager      = "http://localhost:4566"
    ses                 = "http://localhost:4566"
    sns                 = "http://localhost:4566"
    sqs                 = "http://localhost:4566"
    ssm                 = "http://localhost:4566"
    stepfunctions       = "http://localhost:4566"
    sts                 = "http://localhost:4566"
  }
}

locals {
  app_name = "integration-marketo"
  pwd      = "/home/s/Desktop/Projects/Microservice"
}

variable "environment" {
  type = string
  default = "local"
}

############################################################################################################
# SQS
############################################################################################################

# Define the batch queue, used to hold batches defined in the marketo_batch table
resource "aws_sqs_queue" "batch-queue" {
  name                      = "${local.app_name}-batch-queue-${var.environment}"
  delay_seconds             = 1
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 1
  redrive_policy            = jsonencode({
    deadLetterTargetArn     = aws_sqs_queue.batch-queue-deadletter.arn
    maxReceiveCount         = 10
  })
}

# Define the dead letter queue for the batch queue
resource "aws_sqs_queue" "batch-queue-deadletter" {
  name = "${local.app_name}-batch-queue-deadletter-${var.environment}"
}

# Define the policy for the dead letter queue
resource "aws_sqs_queue_redrive_allow_policy" "batch-queue-deadletter" {
  queue_url = aws_sqs_queue.batch-queue-deadletter.id
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.batch-queue.arn]
  })
}

# Define the policy for the batch queue to allow the lambda to send messages to it
resource "aws_sqs_queue_policy" "batch-queue-policy" {
  queue_url = aws_sqs_queue.batch-queue.url
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Action    = "sqs:SendMessage"
        Resource  = aws_sqs_queue.batch-queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_lambda_function.batch-lambda.arn
          }
        }
      }
    ]
  })
}

################################################################################
# Lambdas
################################################################################

# Define the batch lambda
resource "aws_lambda_function" "batch-lambda" {
  function_name = "${local.app_name}-batch-lambda-${var.environment}"
  role          = aws_iam_role.batch-lambda-role.arn
  handler       = "main"
  runtime       = "go1.x"
  # Comment out the following line to turn off hot-reloading
  filename     = "${local.pwd}/bin/Batch/main.zip"
  #s3_bucket     = "hot-reload"
  #s3_key        = "${local.pwd}/lambdas/Batch"
  environment {
    variables   = {
      QUEUE_URL = aws_sqs_queue.batch-queue.url
    }
  }
}

# Define the role for the batch lambda
resource "aws_iam_role" "batch-lambda-role" {
  name = "batch-lambda-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Define the permission for the batch lambda to be invoked by the batch queue
resource "aws_lambda_permission" "batch-lambda-permission" {
  statement_id = "AllowExecutionFromSQS"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.batch-lambda.function_name
  principal = "sqs.amazonaws.com"
  source_arn = aws_sqs_queue.batch-queue.arn
}