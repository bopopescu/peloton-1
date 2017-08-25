import logging
import os
import pytest
import time

from docker import Client
from tools.pcluster.pcluster import setup, teardown


log = logging.getLogger(__name__)


#
# Global setup / teardown across all test suites
#
@pytest.fixture(scope="session", autouse=True)
def setup_cluster(request):
    log.info('setup cluster')
    if os.getenv('CLUSTER', ''):
        log.info('cluster mode')
    else:
        log.info('local pcluster mode')
        setup(enable_peloton=True)
        time.sleep(5)

    def teardown_cluster():
        log.info('\nteardown cluster')
        if os.getenv('CLUSTER', ''):
            log.info('cluster mode, no teardown actions')
        elif os.getenv('NO_TEARDOWN', ''):
            log.info('skip teardown')
        else:
            log.info('teardown, writing logs')
            try:
                cli = Client(base_url='unix://var/run/docker.sock')
                log.info(cli.logs('peloton-jobmgr0'))
                log.info(cli.logs('peloton-jobmgr1'))
            except Exception as e:
                log.info(e)
            teardown()

    request.addfinalizer(teardown_cluster)


class Container(object):
    def __init__(self, names):
        self._cli = Client(base_url='unix://var/run/docker.sock')
        self._names = names

    def start(self):
        for name in self._names:
            self._cli.start(name)
            log.info('%s started', name)

    def stop(self):
        for name in self._names:
            self._cli.stop(name, timeout=0)
            log.info('%s stopped', name)

    def restart(self):
        for name in self._names:
            self._cli.restart(name, timeout=0)
            log.info('%s restarted', name)


@pytest.fixture()
def mesos_master():
    return Container(['peloton-mesos-master'])


@pytest.fixture()
def jobmgr():
    # TODO: We need to pick up the count dynamically.
    return Container(['peloton-jobmgr0', 'peloton-jobmgr1'])


@pytest.fixture()
def mesos_agent():
    # TODO: We need to pick up the count dynamically.
    return Container(['peloton-mesos-agent0', 'peloton-mesos-agent1',
                      'peloton-mesos-agent2'])
