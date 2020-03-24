# Stripe Integration

The Stripe integration is a bridge for the `nCent` core services.

## Getting Started

### Prerequisites

```
GO
Node
npm
Serverless Framework
ginkgo
```

### Installing

Install Serverless Framework

```
npm install -g serverless
```

Install node dependencies

```
npm install
```

## Building

```
make build
```

## Running the tests

```
go test ./...
```

or

```
ginkgo ./...
```

## Setup

1. Create a new `API Key` on Stripe on `https://dashboard.stripe.com/test/apikeys`, this variable is used on `opt:stripe_key` (`stripe_key`)

2. On Stripe dashboard add a new `Product` on `Billing` `https://dashboard.stripe.com/test/subscriptions/products`. Then add a new `Metadata` for the related plan on the `nCent` platform using the variable `PLAN_UUID` and the `uuid` for the related plan. (`default_plan`)

3. After deploy the using `serverless` command go to Stripe on `https://dashboard.stripe.com/test/webhooks` and add the `webhook` using the `/webhook` URL created on `API Gateway`, the event type should be `customer.updated`

4. After create the `webhook` reveal and copy the `Signing Secret` on the same page and paste on the `stripe-integration` function `AWS` as the `WEBHOOK_SECRET`. The idea to make sure that all post really comes from Stripe (`webhook_secret`)

## Deployment

Example of deploying here we should inform the `default_plan`, `stripe_key` and the `webhook_secret`

```bash
# Example of deploy on production
serverless deploy --stage production --verbose  --stripe_key "" --default_plan "" --webhook_secret ""
```

## Built With

* [GO](https://golang.org) - The Language
* [Serverless Framework](https://serverless.com) - Deployment Framework

## Authors

* **Eduardo Nunes Peireira** - *Initial work* - [eduardonunesp](https://gitlab.com/eduardonunesp)
