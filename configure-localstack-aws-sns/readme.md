# Configure localstack to mock AWS services on your local machine
### Here, I've provided an example of creating SNS and subscribe it using localstack CLI.

## 1. Install localStack
docker run -it -d -p 4566:4566 -p 4510-4559:4510-4559 --name aws_localstack localstack/localstack

## 2. List all queues
aws --endpoint-url=http://localhost:4566 sqs list-queues --region=us-west-2
##### Error: Unable to locate credentials. You can configure credentials by running "aws configure".
 
## 3. Configure your profile for localStack
aws configure
AWS Access Key ID [None]: test
AWS Secret Access Key [None]: test
Default region name [None]: us-west-2
Default output format [None]: json

OR

## 3. Login with Admin user
aws --endpoint-url=http://localhost:4566 --profile aws-admin sqs list-queues --region=us-west-2

## 4. Create AWS SNS topic
aws --endpoint-url=http://localhost:4566 sns create-topic --region=us-west-2 --name test-topi
============>Response from localstack
{
    "TopicArn": "arn:aws:sns:us-west-2:000000000000:test-topi"
}

## 5. List Topics
aws --endpoint-url=http://localhost:4566 sns list-topics --region=us-west-2

## 6. Subscribe to SNS topic
aws --endpoint-url=http://localhost:4566 sns subscribe --topic-arn arn:aws:sns:us-west-2:000000000000:test-topi --protocol email --notification-endpoint aartiparmar112@gmail.com
=================>Response from localstack
{
    "SubscriptionArn": "arn:aws:sns:us-west-2:000000000000:test-topi:a03bbf53-5198-41e4-a104-ec1e587040a5"
}

## 7. Publish message to SNS
aws --endpoint-url=http://localhost:4566 sns publish --topic-arn arn:aws:sns:us-west-2:000000000000:test-topi --message "Hello from localStack - Aarti Chhasiya"

## 8.Set SNS Subscription attributes, using the SubscriptionArn from the previous step:
aws --endpoint-url=http://localhost:4566 sns set-subscription-attributes --subscription-arn arn:aws:sns:us-west-2:000000000000:test-topi:a03bbf53-5198-41e4-a104-ec1e587040a5 --attribute-name RawMessageDelivery --attribute-value true

## 9. List subcriptions
aws --endpoint-url=http://localhost:4566 sns list-subscriptions

### You can refer below links for reference
1. https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-completion.html
2. https://docs.localstack.cloud/user-guide/aws/sns/
3. What is ARN- https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
