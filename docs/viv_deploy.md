## viv deploy

Deploy your serverless, this is an alias of chaining commands (upload -> remove -> create)

```
viv deploy src_file[.go|.zip|dir] [flags]
```

### Options

```
      --env stringArray   Set environment variables as key=value
  -h, --help              help for deploy
```

### Options inherited from parent commands

```
      --api string      REST API endpoint (default "https://hosting.vivgrid.com")
      --secret string   app secret
      --tool string     serverless LLM tool name (default "my_first_llm_tool")
```

### SEE ALSO

* [viv](viv.md)	 - Manage your globally deployed Serverless LLM Functions on vivgrid.com from the command line

