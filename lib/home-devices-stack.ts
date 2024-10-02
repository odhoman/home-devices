import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as eventSources from 'aws-cdk-lib/aws-lambda-event-sources';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as kinesis from 'aws-cdk-lib/aws-kinesis';
import { LambdaHelper } from './helper/lambda-helper';
import { ApiGatewayHelper } from './helper/api-gateway-helper';

export class HomeDevicesStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const macHomeIdIndexName = "MacHomeIdIndex"

    // Table
    var homeDevicesTable = this.createHomeDeviceTable(this, "HomeDevices", "id"); 
    this.addGlobalSecondaryIndex(homeDevicesTable, macHomeIdIndexName, "mac", "homeId")

    // Queue
    const homeDevicesQueue = new sqs.Queue(this, 'HomeDevicesSQS', {
      retentionPeriod: cdk.Duration.days(4),
    });

    // Kinesis Stream
    const kinesisStream = new kinesis.Stream(this, 'HomeDevicesKinesisStream', {
      streamMode: kinesis.StreamMode.PROVISIONED,
      shardCount: 2,
    });

    // Lambdas
    const createDeviceLambda = this.createCreateDeviceLambda(homeDevicesTable, macHomeIdIndexName);
    const getDeviceLambda = this.createGetDeviceLambda(homeDevicesTable);
    const updateDeviceLambda = this.createUpdateDeviceLambda(homeDevicesTable);
    const deleteDeviceLambda = this.createDeleteDeviceLambda(homeDevicesTable);
    const kinesisLambda = this.createKinesisLambda(kinesisStream);
    this.createHomeDeviceListenerLambda(this, homeDevicesQueue, homeDevicesTable, macHomeIdIndexName);

    // ApiGateway
    const api = ApiGatewayHelper.createApiGateway(this, 'HomeDevicesApi');
    ApiGatewayHelper.addLambdaIntegration(api, 'v1/device', 'POST', new apigateway.LambdaIntegration(createDeviceLambda));
    ApiGatewayHelper.addLambdaIntegration(api, 'v1/device/{id}', 'GET', new apigateway.LambdaIntegration(getDeviceLambda));
    ApiGatewayHelper.addLambdaIntegration(api, 'v1/device/{id}', 'PUT', new apigateway.LambdaIntegration(updateDeviceLambda));
    ApiGatewayHelper.addLambdaIntegration(api, 'v1/device/{id}', 'DELETE', new apigateway.LambdaIntegration(deleteDeviceLambda));
  }

  private createHomeDeviceTable(scope: Construct, name: string, partitionKeyName: string): dynamodb.Table {
    var homeDevicesTable = new dynamodb.Table(scope, name, {
      partitionKey: { name: partitionKeyName, type: dynamodb.AttributeType.STRING },
      removalPolicy: cdk.RemovalPolicy.RETAIN,
    });

    return homeDevicesTable;
  }

  private addGlobalSecondaryIndex(table: dynamodb.Table, macHomeIdIndexName: string, partitionKey: string, sortKey: string): void {
    table.addGlobalSecondaryIndex({
      indexName: macHomeIdIndexName,
      partitionKey: { name: partitionKey, type: dynamodb.AttributeType.STRING },
      sortKey: { name: sortKey, type: dynamodb.AttributeType.STRING },
      projectionType: dynamodb.ProjectionType.ALL,
    });
  }

  private createHomeDeviceListenerLambda(scope: Construct, homeDevicesQueue: cdk.aws_sqs.Queue, homeDevicesTable: cdk.aws_dynamodb.Table, macHomeIdIndexName: string): lambda.Function {
    var homeDeviceListenerLambda = LambdaHelper.createLambda(scope, 'HomeDeviceListener', 'bootstrap', 'lambdas/cmd/homeDeviceListener', {
      SQS_QUEUE_URL: homeDevicesQueue.queueUrl,
      HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName      
    });

    homeDevicesTable.grantReadWriteData(homeDeviceListenerLambda);

    homeDeviceListenerLambda.addEventSource(new eventSources.SqsEventSource(homeDevicesQueue, {
      batchSize: 10,
    }));

    return homeDeviceListenerLambda;
  }

  private createCreateDeviceLambda(homeDevicesTable: cdk.aws_dynamodb.Table, macHomeIdIndexName: string): cdk.aws_lambda.Function {
    var createDeviceLambda = LambdaHelper.createLambda(this, 'CreateDevice', 'bootstrap', 'lambdas/cmd/createDevice', {
      HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
      MAC_HOMEID_INDEX_NAME: macHomeIdIndexName
    });

    createDeviceLambda.addToRolePolicy(new iam.PolicyStatement({
      actions: ['dynamodb:Query'],
      resources: [
        homeDevicesTable.tableArn,
        `${homeDevicesTable.tableArn}/index/${macHomeIdIndexName}`
      ],
    }));

    homeDevicesTable.grantWriteData(createDeviceLambda);

    return createDeviceLambda;
  }

  private createGetDeviceLambda(homeDevicesTable: cdk.aws_dynamodb.Table): cdk.aws_lambda.Function {
    var getDeviceLambda = LambdaHelper.createLambda(this, 'GetDevice', 'bootstrap', 'lambdas/cmd/getDevice', {
      HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
    });

    homeDevicesTable.grantReadData(getDeviceLambda);

    return getDeviceLambda;
  }

  private createUpdateDeviceLambda(homeDevicesTable: cdk.aws_dynamodb.Table): cdk.aws_lambda.Function {
    var updateDeviceLambda = LambdaHelper.createLambda(this, 'UpdateDevice', 'bootstrap', 'lambdas/cmd/updateDevice', {
      HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
    });

    homeDevicesTable.grantWriteData(updateDeviceLambda);

    return updateDeviceLambda;
  }

  private createDeleteDeviceLambda(homeDevicesTable: cdk.aws_dynamodb.Table): cdk.aws_lambda.Function {
    var deleteDeviceLambda = LambdaHelper.createLambda(this, 'DeleteDevice', 'bootstrap', 'lambdas/cmd/deleteDevice', {
      HOME_DEVICE_TABLE_NAME: homeDevicesTable.tableName,
    });

    homeDevicesTable.grantReadWriteData(deleteDeviceLambda);

    return deleteDeviceLambda;
  }

  private createKinesisLambda(kinesisStream: kinesis.Stream): cdk.aws_lambda.Function {
    var kinesisListener = LambdaHelper.createLambda(this, 'KinesisListener', 'bootstrap', 'lambdas/cmd/kinesisListener', {
      // Puedes agregar tus variables de entorno si las necesitas
    });

    // Agregar Kinesis como Event Source y configurar ParallelizationFactor
    kinesisListener.addEventSource(new eventSources.KinesisEventSource(kinesisStream, {
      startingPosition: lambda.StartingPosition.LATEST,
      batchSize: 10,  // Establecer el batch size a 10
      maxBatchingWindow: cdk.Duration.seconds(10),  // Establecer el batch window a 10 segundos
      parallelizationFactor: 3,
    }));

    return kinesisListener;
  }
}
