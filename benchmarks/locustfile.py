import os
import time
import uuid
import random

import redis

from locust import Locust
from locust import events
from locust import HttpLocust
from locust import TaskSet
from locust import task


class RedisClient(object):
    def __init__(self, host="localhost", port=6379):
        self.rc = redis.StrictRedis(host=host, port=port)

    def execute_command(self, name, command, *args):
        result = None
        start_time = time.time()
        try:
            result = self.rc.execute_command(command, *args)
            if not result:
                result = ''
        except Exception as e:
            total_time = (time.time() - start_time) * 1000
            events.request_failure.fire(
                request_type=command,
                name=name,
                response_time=total_time,
                exception=e
            )
        else:
            total_time = (time.time() - start_time) * 1000
            length = len(str(result).encode('utf-8'))

            events.request_success.fire(
                request_type=command,
                name=name,
                response_time=total_time,
                response_length=length
            )
        return result


class LimiterBehavior(TaskSet):
    def on_start(self):
        self.user_id = uuid.uuid1()
        self.reids = RedisClient(port=46379)

    @task(1)
    def get_tokens_http(self):
        response = self.client.get(
            '/API/v1/limiter/user_{}/1000/1/1000/1/'.format(self.user_id),
            name='/API/v1/limiter/test/1000/1/1000/1/'
        )

        assert response.status_code == 200, response.text

    @task(1)
    def get_tokens_http_(self):
        response = self.client.get(
            '/API/v1/limiter/user_{}_endpoint_{}/1000/1/1000/1/'.format(
                self.user_id, random.randint(1, 1000000)
            ),
            name='/API/v1/limiter/test/1000/1/1000/1/'
        )

        assert response.status_code == 200, response.text

    @task(1)
    def get_tokens_redis(self):
        key = 'user_{}'.format(self.user_id)
        tokens = self.reids.execute_command(
            'REDUCE 1000 1 1000 1',
            'REDUCE',
            key,
            1000,
            1,
            1000,
            1
        )
        assert tokens >= -1, tokens

    @task(1)
    def get_tokens_redis(self):
        key = 'user_{}_endpoint_{}'.format(
            self.user_id, random.randint(1, 1000000)
        )
        tokens = self.reids.execute_command(
            'REDUCE 1000 10 1000 10',
            'REDUCE',
            key,
            1000,
            10,
            1000,
            10
        )
        assert tokens >= -1, tokens


class WebsiteLimiter(HttpLocust):
    task_set = LimiterBehavior
    min_wait = 10
    max_wait = 1000
