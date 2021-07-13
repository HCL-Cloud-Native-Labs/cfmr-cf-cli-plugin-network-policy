# CF CLI Plugin Network Policy 

| Build | State |
| ---     | ---   |
| CI | [![CI status](https://github.com/HCL-Cloud-Native-Labs/cfmr-cf-cli-plugin-network-policy/workflows/Main%20CI/badge.svg)](https://github.com/HCL-Cloud-Native-Labs/cfmr-cf-cli-plugin-network-policy/actions?query=workflow%3AMain+CI) |

## About
This plugin enables App to App direct communication in CFMR

## Usage
Build image from repo:

```shell script
# Clone repo and cd into it then...
docker build --no-cache --tag hclcnlabs/cfmr-cf-cli-plugin-network-policy:1.0.0 .
```

```shell script
# Push the image to dockerhub(hclcnlabs)...
docker push hclcnlabs/cfmr-cf-cli-plugin-network-policy:1.0.0
```
## Installation Steps
- `cf install-plugin plugin -f`

## Maintainers
The following are the custodians of this codebase:

| Name | Email |
| ---     | ---   |
| Manoj | manojkumar.tyagi@hcl.com |
