// Edge cases for testing

// Empty file is tested separately

// Single line
const x = 42;

// Comments only
// This is a comment
/* Multi-line
   comment */

// Nested functions
function outer() {
    function inner() {
        return "nested";
    }
    return inner();
}

// Complex generic constraints
interface Repository<T extends { id: string | number }> {
    findById(id: T["id"]): Promise<T | null>;
    save(entity: T): Promise<T>;
}

// Deeply nested types
type DeepNested = {
    level1: {
        level2: {
            level3: {
                value: string;
            };
        };
    };
};

// Multiple type parameters
function merge<T, U>(obj1: T, obj2: U): T & U {
    return { ...obj1, ...obj2 };
}

// Conditional types
type IsString<T> = T extends string ? true : false;

// Mapped types
type Readonly<T> = {
    readonly [P in keyof T]: T[P];
};

// Template literal types
type EventName = "click" | "focus" | "blur";
type EventHandler = `on${Capitalize<EventName>}`;
