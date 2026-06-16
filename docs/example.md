### Example

You can create a `vivgrid.yml` file in the root of your project to configure the deployment settings.

```yaml
zipper: zipper.vivgrid.com # Port 9000 will be added automatically
# OR specify custom port:
# zipper: zipper.vivgrid.com:8080
secret: <your_app_secret>
tool: <your_function_name>
```

- zipper: Zipper address - can be a domain (port 9000 added automatically) or domain:port (default "zipper.vivgrid.com")
- secret: App secret for authentication
- tool: Serverless LLM Function name (default "my_first_llm_tool")
