package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

func main() {
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Instance Name", "Private IP", "Public IP"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	data := [][]string{}

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			// get the name tag
			instanceName := "Undefined"
			for _, t := range i.Tags {
				if strings.ToLower(aws.StringValue(t.Key)) == "name" {
					instanceName = aws.StringValue(t.Value)
				}
			}

			if instanceName != "Undefined" {
				data = append(data, []string{instanceName, aws.StringValue(i.PrivateIpAddress), aws.StringValue(i.PublicIpAddress)})
			}
		}
	}

	sort.Slice(data, func(i int, j int) bool {
		return data[i][0] < data[j][0]
	})

	table.AppendBulk(data)
	table.Render()

	//fmt.Println(result)
}
