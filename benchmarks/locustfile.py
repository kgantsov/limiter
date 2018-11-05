import os
from random import choice

from locust import HttpLocust
from locust import TaskSet
from locust import task


class UserBehavior(TaskSet):
    def on_start(self):
        pass

    @task(1)
    def get_tokens_http(self):
        response = self.client.get(
            "/API/v1/limiter/test/1000/1/1000/1/",
            name='/API/v1/limiter/test/1000/1/1000/1/'
        )

        assert response.status_code == 200, response.text


class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 1000
    max_wait = 4000
