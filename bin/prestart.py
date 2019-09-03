#!/usr/bin/env python

import os
import sys
from systemd.daemon import notify, Notification

if os.environ.get('SYSTEMD_NOTIFY_ENABLE') == 'true':
    notify(Notification.READY)
    notify(Notification.STATUS, "Service ready")


print("Finished prestart")
sys.exit(0)
