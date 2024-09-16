package fn

import (
    "github.com/therecipe/qt/widgets"
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/gui"
)

// CreateMainWindow 创建和返回主窗口
func Qtgui() *widgets.QMainWindow {
    // 创建主窗口
    window := widgets.NewQMainWindow(nil, 0)
    window.SetWindowTitle("Bat")
    window.SetMinimumSize2(400, 300)
    window.SetWindowIcon(gui.NewQIcon5("Bat.icon"))
    
    // 设置窗口样式
    window.SetStyleSheet(`
        QMainWindow {
            background: qlineargradient(x1:0, y1:0, x2:1, y2:1, stop:0 #f0f0f0, stop:1 #e0e0e0);
        }
        QLabel {
            font-family: 'Courier New';
            font-size: 18px;
            padding: 20px;
        }
        QPushButton {
            font-family: 'Arial';
            font-size: 18px;
            padding: 15px;
            border: 2px solid #007bff;
            border-radius: 10px;
            background-color: #007bff;
            color: white;
        }
        QPushButton:hover {
            background-color: #0056b3;
        }
        QPushButton:pressed {
            background-color: #004080;
        }
    `)

    // 创建中心部件和布局
    centralWidget := widgets.NewQWidget(nil, 0)
    mainLayout := widgets.NewQVBoxLayout()
    mainLayout.SetSpacing(5)

    // 创建带颜色的 ASCII 图标
    asciiArt := `
    <pre style="color: #3B78FF;">  <!-- 设置颜色 -->
    ██████╗  █████╗ ████████╗
    ██╔══██╗██╔══██╗╚══██╔══╝
    ██████╔╝███████║   ██║   
    ██╔══██╗██╔══██║   ██║   
    ██████╔╝██║  ██║   ██║   
    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   
    </pre>`

    label := widgets.NewQLabel2(asciiArt, nil, 0)
    label.SetAlignment(core.Qt__AlignCenter)
    mainLayout.AddWidget(label, 0, core.Qt__AlignCenter)

    // 创建水平布局来并排放置按钮
    buttonLayout := widgets.NewQHBoxLayout()

    // 创建按钮并添加到布局中
    button := widgets.NewQPushButton2("启动", nil)
    exit := widgets.NewQPushButton2("退出", nil)

    // 调整按钮布局的间距和对齐方式
    buttonLayout.SetSpacing(5) // 设置按钮之间的间距为5像素
    buttonLayout.SetContentsMargins(10, 10, 10, 10) // 设置按钮布局的外边距

    buttonLayout.AddWidget(button, 0, core.Qt__AlignCenter)
    buttonLayout.AddWidget(exit, 0, core.Qt__AlignCenter)

    // 将按钮布局添加到主布局中
    mainLayout.AddLayout(buttonLayout, 0)

    // 设置按钮点击事件处理函数
    button.ConnectClicked(func(checked bool) {
        label.SetText("开始扫描")
    })

    exit.ConnectClicked(func(checked bool) {
        window.Close()  // 关闭窗口
    })

    // 设置中心部件的布局
    centralWidget.SetLayout(mainLayout)
    window.SetCentralWidget(centralWidget)

    return window
}
