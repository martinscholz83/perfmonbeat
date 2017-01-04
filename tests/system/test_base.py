from perfmonbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Perfmonbeat normally
        """
        self.render_config_template(
                path=os.path.abspath(self.working_dir) + "/log/*"
        )

        perfmonbeat_proc = self.start_beat()
        self.wait_until( lambda: self.log_contains("perfmonbeat is running"))
        exit_code = perfmonbeat_proc.kill_and_wait()
        assert exit_code == 0
