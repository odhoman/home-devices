import * as cdk from 'aws-cdk-lib';
import { Template, Match } from 'aws-cdk-lib/assertions';
import { HomeDevicesStack } from '../lib/home-devices-stack';

test('DynamoDB Table Created', () => {
    const app = new cdk.App();
    const stack = new HomeDevicesStack(app, 'MyTestStack');
    const template = Template.fromStack(stack);

    template.hasResourceProperties('AWS::DynamoDB::Table', {
        KeySchema: Match.arrayWith([{
            AttributeName: 'id',
            KeyType: 'HASH'
        }]),
        ProvisionedThroughput: {
            ReadCapacityUnits: 5,
            WriteCapacityUnits: 5
        },
        GlobalSecondaryIndexes: Match.arrayWith([
            Match.objectLike({
                IndexName: 'MacHomeIdIndex',
                KeySchema: [
                    {
                        AttributeName: 'mac',
                        KeyType: 'HASH',  // Partition key
                    },
                    {
                        AttributeName: 'homeId',
                        KeyType: 'RANGE', // Sort key
                    },
                ],
                Projection: {
                    ProjectionType: 'ALL',
                },
            }),
        ]),
    });
});

test('SQS Queue Created', () => {
    const app = new cdk.App();
    const stack = new HomeDevicesStack(app, 'MyTestStack');
    const template = Template.fromStack(stack);

    template.hasResourceProperties('AWS::SQS::Queue', {
        MessageRetentionPeriod: 345600 // 4 days in seconds
    });
});

test('Lambda Functions Created', () => {
    const app = new cdk.App();
    const stack = new HomeDevicesStack(app, 'MyTestStack');
    const template = Template.fromStack(stack);

    // Check CreateDevice Lambda
    template.hasResourceProperties('AWS::Lambda::Function', {
        Handler: 'bootstrap',
        Runtime: 'provided.al2023',
        Role: Match.objectLike({
            "Fn::GetAtt": [
                Match.stringLikeRegexp('CreateDeviceServiceRole'),
                "Arn"
            ]
        }),
        Environment: {
            Variables: {
                HOME_DEVICE_TABLE_NAME: Match.anyValue(),
                MAC_HOMEID_INDEX_NAME: Match.anyValue()
            }
        }
    });


    // Check GetDevice Lambda
    template.hasResourceProperties('AWS::Lambda::Function', {
        Handler: 'bootstrap',
        Runtime: 'provided.al2023',
        Role: Match.objectLike({
            "Fn::GetAtt": [
                Match.stringLikeRegexp('GetDeviceServiceRole'),
                "Arn"
            ]
        }),
        Environment: {
            Variables: {
                HOME_DEVICE_TABLE_NAME: Match.anyValue(),
            }
        }
    });

    // Check deleteDevice Lambda
    template.hasResourceProperties('AWS::Lambda::Function', {
        Handler: 'bootstrap',
        Runtime: 'provided.al2023',
        Role: Match.objectLike({
            "Fn::GetAtt": [
                Match.stringLikeRegexp('DeleteDeviceServiceRole'),
                "Arn"
            ]
        }),
        Environment: {
            Variables: {
                HOME_DEVICE_TABLE_NAME: Match.anyValue(),
            }
        }
    });

    // Check updateDevice Lambda
    template.hasResourceProperties('AWS::Lambda::Function', {
        Handler: 'bootstrap',
        Runtime: 'provided.al2023',
        Role: Match.objectLike({
            "Fn::GetAtt": [
                Match.stringLikeRegexp('UpdateDeviceServiceRole'),
                "Arn"
            ]
        }),
        Environment: {
            Variables: {
                HOME_DEVICE_TABLE_NAME: Match.anyValue(),
            }
        }
    });

    // Check homeDeviceListener Lambda
    template.hasResourceProperties('AWS::Lambda::Function', {
        Handler: 'bootstrap',
        Runtime: 'provided.al2023',
        Role: Match.objectLike({
            "Fn::GetAtt": [
                Match.stringLikeRegexp('HomeDeviceListener'),
                "Arn"
            ]
        }),
        Environment: {
            Variables: {
                HOME_DEVICE_TABLE_NAME: Match.anyValue(),
                SQS_QUEUE_URL: Match.anyValue(),
            }
        }
    });
});


test('API Gateway Methods Created', () => {
    const app = new cdk.App();
    const stack = new HomeDevicesStack(app, 'MyTestStack');
    const template = Template.fromStack(stack);
  
    // Verificar que el método POST para el recurso 'v1/device' esté creado
    template.hasResourceProperties('AWS::ApiGateway::Method', {
      HttpMethod: 'POST',
      ResourceId: Match.anyValue(),  // Asegurando que haya un ResourceId asociado al método
      RestApiId: Match.anyValue(),  // Asegurando que haya un RestApiId asociado al método
      Integration: {
        IntegrationHttpMethod: 'POST',
        Type: 'AWS_PROXY',
      },
    });
  
    // Verificar que el método GET para el recurso 'v1/device/{id}' esté creado
    template.hasResourceProperties('AWS::ApiGateway::Method', {
      HttpMethod: 'GET',
      ResourceId: Match.anyValue(),
      RestApiId: Match.anyValue(),
      Integration: {
        IntegrationHttpMethod: 'POST',
        Type: 'AWS_PROXY',
      },
    });
  
    // Verificar que el método PUT para el recurso 'v1/device/{id}' esté creado
    template.hasResourceProperties('AWS::ApiGateway::Method', {
      HttpMethod: 'PUT',
      ResourceId: Match.anyValue(),
      RestApiId: Match.anyValue(),
      Integration: {
        IntegrationHttpMethod: 'POST',
        Type: 'AWS_PROXY',
      },
    });
  
    // Verificar que el método DELETE para el recurso 'v1/device/{id}' esté creado
    template.hasResourceProperties('AWS::ApiGateway::Method', {
      HttpMethod: 'DELETE',
      ResourceId: Match.anyValue(),
      RestApiId: Match.anyValue(),
      Integration: {
        IntegrationHttpMethod: 'POST',
        Type: 'AWS_PROXY',
      },
    });
  });