import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';

export class LambdaHelper {
  public static createLambda(scope: Construct, id: string, handler: string, codePath: string, environment: { [key: string]: string }): lambda.Function {
    return new lambda.Function(scope, id, {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: handler,
      code: lambda.Code.fromAsset(codePath),
      timeout: cdk.Duration.seconds(5),
      environment: environment,
    });
  }
}
