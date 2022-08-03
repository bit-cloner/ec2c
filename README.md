### What is EC2C
EC2-C stands for EC2 Copy : Copy/migrate EC2 instances from one AWS account to another AWS account from the comfort of a commnad line.

### What can EC2C do
1. Shows you a list of EC2 instances from the AWS account and chosen Region as per AWS credentials
2. Creates Amaxzon Machine Images for chosen EC2 instances.
3. Waits untill the images are in "available" status. Time taken depends on time taken for a snapshot to be created. Depends on the size of disk. For timeliness timeout occurs after 40 minutes.
4. Asks for target AWS account number. Thsi is the AWS account where you would want the EC2 instances to be migrated.
5. Changes permissions on newly created images. Adds the target account as a shared account for the image.

### Prerequisites
1. AWS credentials from source account with appropriate persmissions to create an AMI 
2. Target AWS account number

### How to get it
Chose the right artifact for your CPU architecture and OS type from https://github.com/bit-cloner/ec2c/releases
```
wget https://github.com/bit-cloner/ec2c/releases/download/0.9.1/ec2c-0.9.1-linux-amd64.tar.gz
```
```
tar -xvzf ec2c-0.9.1-linux-amd64.tar.gz
```
### How to use it
```
chmod +x ./ec2c
```
```
./ec2c
```