export function SelectWorkspace() { return Promise.resolve("dummy-ws-id"); }
export function ListTasks() { return Promise.resolve([]); }
export function GetWorkspace(arg1) { return Promise.resolve({}); }
export function CreateTask(arg1, arg2) { return Promise.resolve({}); }
export function RunTask(arg1) { return Promise.resolve(); }
export function ListAttempts(_taskId) { return Promise.resolve([]); }
export function GetPoolSummaries() { return Promise.resolve([]); }
