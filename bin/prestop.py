#!/usr/bin/env python

import time
import os
import sys
from systemd.daemon import notify, Notification

if os.environ.get('SYSTEMD_NOTIFY_ENABLE') == 'true':
    notify(Notification.STOPPING)
    notify(Notification.STATUS, "Service stopping...")
time.sleep(3)

print("Finished prestart")
sys.exit(0)
