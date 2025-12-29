// URLDB Plugin System TypeScript Definitions

declare global {
  // 应用接口
  interface App {
  }

  // URL 模型
  interface URL {
    id: string;
    url: string;
    title: string;
    category: string;
    tags: string[];
    createdAt: Date;
    updatedAt: Date;
  }

  // 用户模型
  interface User {
    id: string;
    username: string;
    email: string;
    createdAt: Date;
  }

  // 钩子事件
  interface URLEvent {
    app: App;
    url: URL;
    data: Record<string, any>;
    next(): void;
  }

  interface UserEvent {
    app: App;
    user: User;
    data: Record<string, any>;
    next(): void;
  }

  interface APIEvent {
    app: App;
    request: any;
    path: string;
    method: string;
    headers: Record<string, string>;
    body: any;
    next(): void;
  }
}

// 钩子函数声明
declare function onURLAdd(handler: (e: URLEvent) => void): void;
declare function onURLAccess(handler: (e: URLAccessEvent) => void): void;
declare function onUserLogin(handler: (e: UserEvent) => void): void;

// 路由函数声明
declare function routerAdd(method: string, path: string, handler: (ctx: any) => void): void;

// 定时任务函数声明
declare function cronAdd(name: string, schedule: string, handler: () => void): void;

// 全局变量
declare const $app: App;
declare const __hooks: string;

export {};
