# DataDog Auto-Discovery Annotations

If you are planning to monitor your Pega Helm Chart deployment using the DataDog Helm Chart or DaemonSet, you can set the following value in your values.yaml file:
```
global:
  enableDataDogAutoDiscovery: true
```

Proving this setting will cause helm to add annotations to the Pega containers that enable DataDog to discover, connect and collect metrics via JMX.


# Metric Selection

To add, limit or change which metrics are gathered by DataDog, modify the pega-web-tomcat-metrics.yaml file contained in this directory.

See [https://docs.datadoghq.com/integrations/java/?tab=host] for more information on how to select which metrics to collect.

You should only modify the init_config.conf section as the other configuration is specific to how DataDog connects to deployed Pega containers.


# DataDog Deployment Considerations

When deploying the DataDog Helm Chart, you'll need to make sure that it uses the JMX enabled DataDog agent container.

You can do this by specifying this in the values.yaml file used for the DataDog Helm Chart:
```
agents:
  image: tag=7-jmx
```


# DataDog Annotations

The generates annotations will resemble:
```
kind: Pod
metadata:
  annotations:
    ad.datadoghq.com/pega-web-tomcat.check_names: '["pega-web-tomcat-c27ca"]'
    ad.datadoghq.com/pega-web-tomcat.init_configs: '[{"collect_default_metrics":false,"conf":[{"include":{"attribute":{"currentThreadCount":{"alias":"tomcat.threads.count","metric_type":"gauge"},"currentThreadsBusy":{"alias":"tomcat.threads.busy","metric_type":"gauge"},"maxThreads":{"alias":"tomcat.threads.max","metric_type":"gauge"}},"type":"ThreadPool"}},{"include":{"attribute":{"bytesReceived":{"alias":"tomcat.bytes_rcvd","metric_type":"counter"},"bytesSent":{"alias":"tomcat.bytes_sent","metric_type":"counter"},"errorCount":{"alias":"tomcat.error_count","metric_type":"counter"},"maxTime":{"alias":"tomcat.max_time","metric_type":"gauge"},"processingTime":{"alias":"tomcat.processing_time","metric_type":"counter"},"requestCount":{"alias":"tomcat.request_count","metric_type":"counter"}},"type":"GlobalRequestProcessor"}},{"include":{"attribute":{"errorCount":{"alias":"tomcat.servlet.error_count","metric_type":"counter"},"processingTime":{"alias":"tomcat.servlet.processing_time","metric_type":"counter"},"requestCount":{"alias":"tomcat.servlet.request_count","metric_type":"counter"}},"j2eeType":"Servlet"}},{"include":{"attribute":{"accessCount":{"alias":"tomcat.cache.access_count","metric_type":"counter"},"hitsCounts":{"alias":"tomcat.cache.hits_count","metric_type":"counter"}},"type":"Cache"}},{"include":{"attribute":{"jspCount":{"alias":"tomcat.jsp.count","metric_type":"counter"},"jspReloadCount":{"alias":"tomcat.jsp.reload_count","metric_type":"counter"}},"type":"JspMonitor"}}],"is_jmx":true}]'
    ad.datadoghq.com/pega-web-tomcat.instances: '[{"host":"%%host%%","port":9001}]'
```


# Troubleshooting

To confirm DataDog monitoring is working correctly, you can exec into the DataDog agent pod:
```
kubectl -n kube-system exec -it datadog-monitoring-jtmdk /bin/bash

root@datadog-monitoring-jtmdk:/# agent status
Getting the status from the agent.

...

========
JMXFetch
========

  Initialized checks
  ==================
    pega-web-tomcat-c27ca
      instance_name : pega-web-tomcat-c27ca-10.0.19.158-9001
      message : <no value>
      metric_count : 175
      service_check_count : 0
      status : OK
    pega-web-tomcat-f3056
      instance_name : pega-web-tomcat-f3056-10.0.16.14-9001
      message : <no value>
      metric_count : 175
      service_check_count : 0
      status : OK
  Failed checks
  =============
    no checks

...

root@datadog-monitoring-jtmdk:/#
```

You can additionally use the ```agent jmx list collected``` and the ```agent jmx list everything``` for more info.

