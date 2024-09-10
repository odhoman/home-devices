import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayHelper {
  public static createApiGateway(scope: Construct, apiName: string): apigateway.RestApi {
    return new apigateway.RestApi(scope, apiName, {
      restApiName: apiName,
      description: 'This service serves as the API for managing home devices.',
    });
  }

  public static addLambdaIntegration(api: apigateway.RestApi, path: string, method: string, lambdaIntegration: apigateway.LambdaIntegration) {
    const resource = api.root.resourceForPath(path);
    resource.addMethod(method, lambdaIntegration);
  }
}
