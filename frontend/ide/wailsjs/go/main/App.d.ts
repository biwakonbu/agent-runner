export function SelectWorkspace():Promise<string>;
export function ListTasks():Promise<Array<any>>;
export function GetWorkspace(arg1:string):Promise<any>;
export function CreateTask(arg1:string, arg2:string):Promise<any>;
export function RunTask(arg1:string):Promise<void>;
