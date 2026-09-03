<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>500 · 服务器内部错误</title>
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
            background: #f0f4f9;
            padding: 16px;
        }
        .card {
            max-width: 360px;
            width: 100%;
            background: white;
            border-radius: 24px;
            padding: 28px 22px 22px;
            box-shadow: 0 8px 24px rgba(0,0,0,0.06);
            text-align: center;
        }

        /* 图标 + 标题 */
        .icon { font-size: 48px; line-height: 1; margin-bottom: 4px; }
        .title { font-size: 20px; font-weight: 600; color: #111; margin-bottom: 2px; }
        .desc { font-size: 14px; color: #6a7a8e; margin-bottom: 16px; }

        /* 详情信息 */
        .detail {
            background: #f6f9fe;
            border-radius: 14px;
            padding: 12px 14px;
            font-size: 14px;
            color: #1f2a3a;
            margin-bottom: 18px;
            border: 1px solid #e6edf5;
            text-align: left;
        }
        .detail span { color: #7a8ca0; }


    </style>
</head>
<body>

<div class="card">
    <div class="icon" id="icon">⚠️</div>
    <div class="title" id="title">500 · 服务器内部错误</div>
    <div class="desc" id="desc">服务器遇到意外情况，无法完成请求。</div>
    <div class="detail" id="detail">
        <span>{{ .error }}</span>
    </div>
</div>
</body>
</html>
