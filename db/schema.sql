-- USERS Table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,      
    email TEXT UNIQUE NOT NULL,                
    username TEXT UNIQUE NOT NULL,             
    password TEXT NOT NULL,              
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP 
);

-- POSTS Table
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,     
    user_id INTEGER NOT NULL,                 
    title TEXT NOT NULL,                      
    content TEXT NOT NULL,                    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) 
);

-- COMMENTS Table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,     
    post_id INTEGER NOT NULL,                 
    user_id INTEGER NOT NULL,                 
    content TEXT NOT NULL,                    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- CATEGORIES Table
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,     
    name TEXT UNIQUE NOT NULL,                
    description TEXT                          
);




