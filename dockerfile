FROM centos:7

# 维护者信息
LABEL maintainer="uzdz"
LABEL maintainer_email="devmen@163.com"

RUN mkdir /soft
COPY ["rooms_robot", "/soft"]

# 默认启动参数
ENV NOTICE_CHANNEL=
ENV NOTICE_CHANNEL_URL=
ENV NOTICE_CHANNEL_KEY=
ENV PROXY_URL=
ENV TASK_INTERVAL=300

# 添加一个脚本
COPY ["entrypoint.sh", "/soft/entrypoint.sh"]
RUN chmod +x /soft/entrypoint.sh

# 使用脚本作为启动命令
ENTRYPOINT ["/soft/entrypoint.sh"]