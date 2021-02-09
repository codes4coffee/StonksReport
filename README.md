# Stonks Report

![Image of Stonks Meme Man](https://compote.slate.com/images/926e5009-c10a-48fe-b90e-fa0760f82fcd.png?width=1200&rect=680x453&offset=0x30)

Wanted the opportunity to do some work in AWS. This Go app uses the Polygon API to get GMEs high and low from the last trading day, and sends the update to an SNS topic. The app is triggered through CloudWatch and AWS Batch.