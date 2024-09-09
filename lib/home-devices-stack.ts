import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as eventSources from 'aws-cdk-lib/aws-lambda-event-sources';
import * as iam from 'aws-cdk-lib/aws-iam';

export class HomeDevicesStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

        const homeDevicesTable = new dynamodb.Table(this, 'HomeDevices', {
          partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
          removalPolicy: cdk.RemovalPolicy.RETAIN, 
        });

        const macHomeIdIndexName = "MacHomeIdIndex"

        homeDevicesTable.addGlobalSecondaryIndex({
          indexName: macHomeIdIndexName,
          partitionKey: { name: 'mac', type: dynamodb.AttributeType.STRING },
          sortKey: { name: 'homeId', type: dynamodb.AttributeType.STRING },
          projectionType: dynamodb.ProjectionType.ALL,
        });
    
        const homeDevicesQueue = new sqs.Queue(this, 'HomeDevicesSQS', {
          retentionPeriod: cdk.Duration.days(4),
        });
    
        // Lambda GetDevice
        const getDevice = new lambda.Function(this, 'GetDevice', {
          runtime: lambda.Runtime.PROVIDED_AL2023,
          handler: 'bootstrap',
          code: lambda.Code.fromAsset('lambdas/getDevice'),
          timeout: cdk.Duration.seconds(5),
          environment: {
            HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
            MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
          },
        });
    
        // Lambda UpdateDevice
        const updateDevice = new lambda.Function(this, 'UpdateDevice', {
          runtime: lambda.Runtime.PROVIDED_AL2023,
          handler: 'bootstrap',
          code: lambda.Code.fromAsset('lambdas/updateDevice'),
          timeout: cdk.Duration.seconds(5),
          environment: {
            HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
            MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
          },
        });
    
        // Lambda CreateDevice
        const createDevice = new lambda.Function(this, 'CreateDevice', {
          runtime: lambda.Runtime.PROVIDED_AL2023,
          handler: 'bootstrap',
          code: lambda.Code.fromAsset('lambdas/createDevice'),
          timeout: cdk.Duration.seconds(5),
          environment: {
            HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
            MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
          },
        });
    
        // Lambda DeleteDevice
        const deleteDevice = new lambda.Function(this, 'DeleteDevice', {
          runtime: lambda.Runtime.PROVIDED_AL2023,
          handler: 'bootstrap',
          code: lambda.Code.fromAsset('lambdas/deleteDevice'),
          timeout: cdk.Duration.seconds(5),
          environment: {
            HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
            MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
          },
        });
    
        // Lambda HomeDeviceListener
        const homeDeviceListener = new lambda.Function(this, 'HomeDeviceListener', {
          runtime: lambda.Runtime.PROVIDED_AL2023,
          handler: 'bootstrap',
          code: lambda.Code.fromAsset('lambdas/homeDeviceListener'),
          timeout: cdk.Duration.seconds(5),
          environment: {
            SQS_QUEUE_URL: homeDevicesQueue.queueUrl,
            HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
            MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
          },
        });
    
        homeDeviceListener.addEventSource(new eventSources.SqsEventSource(homeDevicesQueue, {
          batchSize: 10,
        }));
    
        // Asignar política IAM para permitir que la Lambda CreateDevice realice consultas en el índice
        createDevice.addToRolePolicy(new iam.PolicyStatement({
          actions: ['dynamodb:Query'],
          resources: [
            homeDevicesTable.tableArn,
            `${homeDevicesTable.tableArn}/index/${macHomeIdIndexName}`
          ],
        }));
    
        // DB Access
        homeDevicesTable.grantReadData(getDevice);
        homeDevicesTable.grantWriteData(updateDevice);
        homeDevicesTable.grantWriteData(createDevice);
        homeDevicesTable.grantReadWriteData(deleteDevice);
        homeDevicesTable.grantReadWriteData(homeDeviceListener);
    
        // Define API Gateway REST API
        const api = new apigateway.RestApi(this, 'HomeDevicesApi', {
          restApiName: 'Home Devices Service',
          description: 'This service serves as the API for managing home devices.',
        });
    
        // /v1 resource
        const v1 = api.root.addResource('v1');
        const device = v1.addResource('device');
        const singleDevice = device.addResource('{id}');
    
        // /v1/device resource for POST (CreateDevice)
        const createDeviceIntegration = new apigateway.LambdaIntegration(createDevice);
        device.addMethod('POST', createDeviceIntegration);
    
        const getDeviceIntegration = new apigateway.LambdaIntegration(getDevice);
        singleDevice.addMethod('GET', getDeviceIntegration);
    
        const updateDeviceIntegration = new apigateway.LambdaIntegration(updateDevice);
        singleDevice.addMethod('PUT', updateDeviceIntegration);
    
        const deleteDeviceIntegration = new apigateway.LambdaIntegration(deleteDevice);
        singleDevice.addMethod('DELETE', deleteDeviceIntegration);

/*     const homeDevicesTable = dynamodb.Table.fromTableName(this, 'HomeDevices', 'HomeDevices');
    const homeDevicesQueue = sqs.Queue.fromQueueArn(this, 'HomeDevicesSQS', 'arn:aws:sqs:us-east-1:798361843197:HomeDevicesSQS');

    const getDevice = new lambda.Function(this, 'GetDevice', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambdas/getDevice'),
      timeout: cdk.Duration.seconds(5),
      environment: {
        HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
      },
    });

    const updateDevice = new lambda.Function(this, 'UpdateDevice', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambdas/updateDevice'),
      timeout: cdk.Duration.seconds(5),
      environment: {
        HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
      },
    });

    const createDevice = new lambda.Function(this, 'CreateDevice', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambdas/createDevice'),
      timeout: cdk.Duration.seconds(5),
      environment: {
        HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
      },
    });

    const deleteDevice = new lambda.Function(this, 'DeleteDevice', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambdas/deleteDevice'),
      timeout: cdk.Duration.seconds(5),
      environment: {
        HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
      },
    });

    const homeDeviceListener = new lambda.Function(this, 'HomeDeviceListener', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambdas/homeDeviceListener'),
      timeout: cdk.Duration.seconds(5),
      environment: {
        SQS_QUEUE_URL: homeDevicesQueue.queueUrl,
        GET_DEVICE_ARN: getDevice.functionArn,
        UPDATE_DEVICE_ARN: updateDevice.functionArn,
      },
    });

    homeDeviceListener.addEventSource(new eventSources.SqsEventSource(homeDevicesQueue, {
      batchSize: 10,
    }));


    createDevice.addToRolePolicy(new iam.PolicyStatement({
      actions: ['dynamodb:Query'],
      resources: [
        `arn:aws:dynamodb:us-east-1:${this.account}:table/HomeDevices/index/MacHomeIdIndex`
      ],
    }));

    // DB Access
    homeDevicesTable.grantReadData(getDevice);
    homeDevicesTable.grantWriteData(updateDevice);
    homeDevicesTable.grantWriteData(createDevice);
    homeDevicesTable.grantReadWriteData(deleteDevice);
    homeDevicesTable.grantReadWriteData(homeDeviceListener);

    // Define API Gateway REST API
    const api = new apigateway.RestApi(this, 'HomeDevicesApi', {
      restApiName: 'Home Devices Service',
      description: 'This service serves as the API for managing home devices.',
    });

    // /v1 resource
    const v1 = api.root.addResource('v1');
    const device = v1.addResource('device');
    const singleDevice = device.addResource('{id}');

    // /v1/device resource for POST (CreateDevice)
    const createDeviceIntegration = new apigateway.LambdaIntegration(createDevice);
    device.addMethod('POST', createDeviceIntegration);


    const getDeviceIntegration = new apigateway.LambdaIntegration(getDevice);
    singleDevice.addMethod('GET', getDeviceIntegration);

    const updateDeviceIntegration = new apigateway.LambdaIntegration(updateDevice);
    singleDevice.addMethod('PUT', updateDeviceIntegration);


    const deleteDeviceIntegration = new apigateway.LambdaIntegration(deleteDevice);
    singleDevice.addMethod('DELETE', deleteDeviceIntegration); */

  }


}
