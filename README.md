# pdf2png-demo

A quick demo of how to use container images with AWS Lambda.

Read the [blog post](https://medium.com/@hichaelmart/using-container-images-with-aws-lambda-7ffbd23697f1) to follow how this works.

## Working locally

1. Authenticate CodeBuild with GitHub

2. Launch the infra stack:

```bash
aws cloudformation deploy \
  --stack-name pdf2png-infra \
  --template-file infra/template.yml \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides ProjectRepository=https://github.com/youruser/pdf2png-demo
```

3. Tag and push and watch your app deploy!

```bash
git tag v1
git push --tags
```

## Testing

Local testing:

```bash
make test
```

Local integration testing:

```bash
make integration-test
```

And then in another terminal:

```bash
BUCKET=mybucket KEY=my/input.pdf make test-event
```
