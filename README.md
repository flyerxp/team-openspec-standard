最终规则效果
Kitex 项目永远不会生成 Controller
Hertz 项目永远不会生成 RPC Handler
DAL、logic、编码风格全团队统一
3. 统一 AI 代码生成口径
以前：AI 有时候给你乱分层、乱建目录、混用 Web/RPC 结构。
现在：全局 AGENTS.md 最高优先级AI 强制：
Kitex 必须 Handler→Service→DAL
Hertz 必须 Controller→Service→DAL
DAL 必须按 {db_name}/gorml/where 分库
所有结构体必须聚合初始化
所有 JSON 必须反引号
禁止变量遮蔽
全部中文注释
所有项目 AI 生成代码 100% 统一风格
4. 统一团队代码评审标准
