FROM public.ecr.aws/lambda/go:1

# Copy function code
COPY main ${LAMBDA_TASK_ROOT}

# Set the CMD to your handler
CMD [ "main" ]