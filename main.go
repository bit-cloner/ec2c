package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

// create a function to create image from instance
/*func CreateImage(c context.Context, instanceID []string, targetaccount []string, awsclient *ec2.NewFromConfig())(imageID []string, err error) {
	// Create image of selected instance

		return imageID, err

}
// create a function to delete image from source account

*/

func main() {
	regions := []string{
		"ap-south-1",
		"eu-west-3",
		"eu-north-1",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-northeast-1",
		"sa-east-1",
		"ca-central-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"cn-north-1",
		"cn-northwest-1",
	}
	region := ""
	prompt := &survey.Select{
		Message: "Select one of the follwing AWS Regions",
		Options: regions,
	}
	survey.AskOne(prompt, &region, nil)
	fmt.Println(region)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}
	fmt.Println("Instance IDs:")
	var instanceIds []string
	for _, r := range result.Reservations {
		//fmt.Println("Reservation ID: " + *r.ReservationId)

		for _, i := range r.Instances {
			//fmt.Println("   " + *i.InstanceId)
			instanceIds = append(instanceIds, *i.InstanceId)
		}

	}
	fmt.Println(instanceIds)

	var multiQs = []*survey.Question{
		{
			Name: "Instances",
			Prompt: &survey.MultiSelect{
				Message: "Choose instances to copy / migrate to another AWS account :",
				Options: instanceIds,
			},
		},
	}
	var singleInput = []*survey.Question{
		{
			Name: "Target Account",
			Prompt: &survey.Input{
				Message: "Enter the target account ID:",
			},
		},
	}
	selectedinstances := []string{}
	err = survey.Ask(multiQs, &selectedinstances)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("you chose: %s\n", strings.Join(selectedinstances, ", "))
	//ask for input of destination account
	targetaccount := ""
	err = survey.Ask(singleInput, &targetaccount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Create images of selected instances
	for _, instanceId := range selectedinstances {
		timestamp := time.Now().Format("2006-01-02-15-04-05")
		namevalue := timestamp + instanceId
		input := &ec2.CreateImageInput{
			InstanceId: &instanceId,
			Name:       &namevalue,
		}
		result, err := client.CreateImage(context.TODO(), input)
		if err != nil {
			fmt.Println("Got an error creating an image of your Amazon EC2 instance:" + instanceId)
			fmt.Println(err)
			return
		}
		fmt.Println("Image created with Image ID: " + *result.ImageId)
		//wait untill image status is available
		describeimageinput := &ec2.DescribeImagesInput{
			ImageIds: []string{*result.ImageId},
		}
		imagestatus := "pending"
		count := 1
		imgCreationTimedout := false
		for imagestatus == "pending" {
			time.Sleep(time.Second * 10)
			count++
			if count > 240 { //Make this value to 240 which equals to 40 minutes
				fmt.Println("Image creation timed out. Aborting...")
				imgCreationTimedout = true
				break
			}
			result, err := client.DescribeImages(context.TODO(), describeimageinput)
			if err != nil {
				fmt.Println("Got an error retrieving information about image:" + *result.Images[0].ImageId)
				fmt.Println(err)
				return
			}
			if result.Images[0].State == "available" {
				imagestatus = "available"
				break
			}
			fmt.Println("status is " + result.Images[0].State)
			fmt.Println("waiting for image status to be available...")
		}

		if imgCreationTimedout {
			fmt.Println("Aborting creation of image for instance: " + instanceId)
			continue
		}
		//modify image attribute
		modifyinput := &ec2.ModifyImageAttributeInput{
			ImageId: result.ImageId,
			LaunchPermission: &types.LaunchPermissionModifications{
				Add: []types.LaunchPermission{
					{
						UserId: &targetaccount,
					},
				},
			},
		}
		_, err = client.ModifyImageAttribute(context.TODO(), modifyinput)
		if err != nil {
			fmt.Println("Got an error modifying the image attribute of AMI created for Instance:" + instanceId + "Image ID:" + *result.ImageId)
			fmt.Println(err)
			return
		}
		fmt.Println("Image attribute modified successfully for AMI created for Instance:" + instanceId + "Image ID:" + *result.ImageId)

	}
}
