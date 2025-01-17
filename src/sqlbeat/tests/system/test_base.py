from sqlbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Sqlbeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        sqlbeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("sqlbeat is running"))
        exit_code = sqlbeat_proc.kill_and_wait()
        assert exit_code == 0
