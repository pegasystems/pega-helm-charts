# BackingServices Helm chart

The backingservices chart installs services like 'Search and Reporting Service' ( abbreviated as SRS) than can be configured with one or more Pega deployments. 
These backing services can be deployed in their own namespace can can be shared across multiple Pega Infinity environments.

**Example:**
_Single backing service shared across all pega environments:_
backingservice 'Search and Reporting Service' deployed and the service endpoint configured across dev, staging and production pega environments. The service provides isolation of data in a shared setup.

_Multiple backing service deployments:_
You can deploy more than one instance of backing service deployments, in case you want to host a seperate deployment of 'Search and Reporting Service' for non-prod and production pega infinity environments. You need to configure the appropriate service endpoint with the pega infinity deployment values.
