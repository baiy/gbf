<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <base href="{{ .baseUrl }}">
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>自动跳转中....</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: system-ui, -apple-system, sans-serif;
        }

        body {
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 16px;
        }

        .card {
            max-width: 360px;
            width: 100%;
            background: white;
            border-radius: 24px;
            padding: 28px 22px 22px;
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
            text-align: center;
        }

        /* 图标 + 标题 */
        .icon {
            font-size: 40px;
            line-height: 1;
            margin-bottom: 4px;
        }

        .title {
            font-size: 20px;
            font-weight: 600;
            color: #111;
            margin-bottom: 2px;
        }

        .desc {
            font-size: 14px;
            color: #6a7a8e;
            margin-bottom: 16px;
        }

        /* 倒计时 */
        .countdown-row {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;
            margin-bottom: 18px;
        }

        .countdown-row .bar-wrap {
            flex: 1;
            height: 6px;
            background: #eaeff4;
            border-radius: 6px;
            overflow: hidden;
        }

        .countdown-row .bar-wrap .fill {
            height: 100%;
            width: 100%;
            background: #0071e3;
            border-radius: 6px;
            transition: width 0.3s linear;
        }

        .countdown-row .num {
            font-size: 16px;
            font-weight: 600;
            color: #0071e3;
            min-width: 28px;
            text-align: center;
        }

        /* 按钮 */
        .actions {
            display: flex;
            gap: 10px;
        }

        .btn {
            flex: 1;
            padding: 12px 10px;
            border-radius: 40px;
            font-size: 15px;
            font-weight: 500;
            cursor: pointer;
            transition: 0.15s;
            background: #f0f4fa;
            color: #1f2a3a;
            border: 1px solid #dce4ed;
        }

        .btn-primary {
            background: #0071e3;
            color: white;
            border-color: #0071e3;
        }

        .btn-primary:active {
            transform: scale(0.96);
            background: #005bb5;
        }

        .btn-secondary:active {
            transform: scale(0.96);
            background: #e4ebf5;
        }

        .btn-error {
            background: #d43d3d;
            color: white;
            border-color: #d43d3d;
        }

        .btn-error:active {
            transform: scale(0.96);
            background: #b82e2e;
        }

        .btn-success {
            background: #0f7b4b;
            color: white;
            border-color: #0f7b4b;
        }

        .btn-success:active {
            transform: scale(0.96);
            background: #0b623a;
        }
    </style>
</head>
<body>

<div class="card">
    <div class="icon" id="icon"></div>
    <div class="title" id="title"></div>
    <div class="desc" id="desc"></div>
    <div class="countdown-row">
        <span>⏱</span>
        <div class="bar-wrap">
            <div class="fill" id="progressFill" style="width:100%"></div>
        </div>
        <span class="num" id="secondsNum"></span>
        <span style="font-size:13px; color:#7a8ca0;">秒</span>
    </div>
    <div class="actions">
        <button class="btn btn-secondary" id="btnBack">返回</button>
        <button class="btn btn-primary" id="btnGo">立即跳转</button>
    </div>
</div>

<script>
    (function () {
        const type = "{{ .type }}"
        const url = "{{ .url }}"
        const info = "{{ .info }}"

        const jumpTime = "success" ? 3 : 10

        const scene = {
            icon: type === "success" ? '✅' : '❌',
            title: type === "success" ? '操作成功' : "操作失败",
            desc: `${info}`,
            btnBack: '返回',
            btnGo: url !== "" ? '立即跳转' : "",
            goClass: `btn-${type}`,
            backClass: 'btn-secondary',
            jumpUrl: `${url}`
        };
        const btnBack = document.getElementById('btnBack');
        const btnGo = document.getElementById('btnGo');
        const secondsNum = document.getElementById('secondsNum');
        const progressFill = document.getElementById('progressFill');

        let countdown = jumpTime;
        let timer = null;
        let isJumping = false;

        function startCountdown() {
            if (timer) {
                clearInterval(timer);
                timer = null;
            }
            countdown = jumpTime;
            updateDisplay();

            timer = setInterval(() => {
                countdown -= 1;
                updateDisplay();
                if (countdown <= 0) {
                    clearInterval(timer);
                    timer = null;
                    doJump();
                }
            }, 1000);
        }

        function updateDisplay() {
            const sec = Math.max(0, countdown);
            secondsNum.textContent = `${sec}`;
            const percent = (sec / jumpTime) * 100;
            progressFill.style.width = percent + '%';
        }

        function stopCountdown() {
            if (timer) {
                clearInterval(timer);
                timer = null;
            }
        }

        function doJump() {
            if (isJumping) return;
            isJumping = true;
            stopCountdown();
            if (scene.jumpUrl === "") {
                location.href = document.referrer
            } else {
                window.location.href = scene.jumpUrl
            }
        }

        btnGo.addEventListener('click', function () {
            if (isJumping) return;
            stopCountdown();
            doJump();
        });

        btnBack.addEventListener('click', function () {
            if (isJumping) return;
            stopCountdown();
            location.href = document.referrer
        });

        document.getElementById('icon').textContent = scene.icon;
        document.getElementById('title').textContent = scene.title;

        if (scene.desc === "") {
            document.getElementById('desc').style.display = 'none';
        } else {
            document.getElementById('desc').style.display = null;
            document.getElementById('desc').textContent = scene.desc;
        }

        btnGo.className = 'btn ' + scene.goClass;
        btnBack.className = 'btn ' + scene.backClass;
        btnBack.textContent = scene.btnBack;
        if (url === "") {
            btnGo.style.display = 'none';
        } else {
            btnGo.style.display = null;
            btnGo.textContent = scene.btnGo;
        }
        updateDisplay();
        startCountdown();
    })();
</script>
</body>
</html>
