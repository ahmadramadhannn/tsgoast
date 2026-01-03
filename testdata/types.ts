// Type and interface test cases

// Simple interface
interface Person {
    name: string;
    age: number;
}

// Interface with optional properties
interface User {
    id: number;
    username: string;
    email?: string;
    phone?: string;
}

// Interface with readonly properties
interface Config {
    readonly apiKey: string;
    readonly endpoint: string;
    timeout: number;
}

// Interface with methods
interface Calculator {
    add(a: number, b: number): number;
    subtract(a: number, b: number): number;
}

// Interface extending another interface
interface Employee extends Person {
    employeeId: string;
    department: string;
}

// Generic interface
interface Container<T> {
    value: T;
    getValue(): T;
}

// Type alias for primitive
type ID = string | number;

// Type alias for object
type Point = {
    x: number;
    y: number;
};

// Type alias for function
type BinaryOp = (a: number, b: number) => number;

// Union type
type Status = "pending" | "approved" | "rejected";

// Intersection type
type Named = { name: string };
type Aged = { age: number };
type NamedAndAged = Named & Aged;

// Generic type alias
type Result<T, E = Error> =
    | { success: true; value: T }
    | { success: false; error: E };

// Exported types
export interface PublicAPI {
    version: string;
    methods: string[];
}

export type Callback<T> = (data: T) => void;
