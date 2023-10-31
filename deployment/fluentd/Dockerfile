# fluentd/Dockerfile
FROM fluent/fluentd:edge
USER root
RUN gem install elasticsearch -v 8.10
RUN gem install fluent-plugin-elasticsearch
RUN gem install faraday_middleware
RUN gem install fluent-plugin-rewrite-tag-filter
RUN gem install fluent-plugin-multi-format-parser
USER fluent