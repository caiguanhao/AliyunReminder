AliyunReminder
==============

Reminds you when your ECS expires.

BUILD
-----

Run `./build.sh` to build the app.

Or `BUILD_DOCKER=1 ./build.sh` to build it in a Docker container and build a Docker image.

How to get the value of Aliyun cookie `login_aliyunid_ticket`:

1. Open [home.console.aliyun.com](https://home.console.aliyun.com) in a Chrome incognito window.
2. Log in with your Aliyun account.
3. Right click the page and select Inspect Element.
4. Select Resouces and then select Cookies > home.console.aliyun.com on the left.
5. Click right pane and press Ctrl-A/Cmd-A to select all text.
6. Paste the text on your text editor.
7. Locate the line starts with `login_aliyunid_ticket` and then copy the text in second column.

How to get Flowdock token:

1. Log in to Flowdock.
2. Select Inbox Settings of your flow.
3. Copy the Flow API token.

------------

LICENSE: MIT
