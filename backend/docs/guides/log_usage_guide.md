## Logger Usage Guide

This document provides instructions on how to use the global logger in your project to log normal events and error logs effectively.

---

### **Initialization**
The logger is globally initialized in the `logger.go` file. It uses the `zap` logging library and stores logs in MongoDB collections for `logs` and `error_logs`.

To use the logger, import it into your module:
```go
import "backend/utils/logger"
```

---

### **Logging Normal Events**
For general logs, use the `Info` method:
```go
logger.Info("File uploaded successfully", map[string]interface{}{
    "user_id": "12345",
    "file_name": "example.txt",
    "action": "upload",
})
```

- **Message**: The first argument is a string describing the log event.
- **Fields**: The second argument is a `map[string]interface{}` containing key-value pairs for additional log details.

---

### **Logging Errors**
For error logs, use the `Error` method:
```go
logger.Error("Failed to upload file", map[string]interface{}{
    "user_id": "12345",
    "file_name": "example.txt",
    "action": "upload",
    "error": err.Error(),
})
```

- **Message**: The first argument is a string describing the error event.
- **Fields**: The second argument contains key-value pairs for context about the error, including the actual error message (`err.Error()`).

---

### **Example in a Service File**
Here’s an example of how to use the logger in a file upload service:
```go
package service

import (
    "backend/utils/logger"
    "errors"
)

func UploadFile(userID, fileName string) error {
    // Simulate an upload action
    success := false // Change this based on actual implementation

    if success {
        logger.Info("File uploaded successfully", map[string]interface{}{
            "user_id": userID,
            "file_name": fileName,
            "action": "upload",
        })
        return nil
    }

    err := errors.New("upload failed due to network issue")
    logger.Error("Failed to upload file", map[string]interface{}{
        "user_id": userID,
        "file_name": fileName,
        "action": "upload",
        "error": err.Error(),
    })

    return err
}
```

---

### **Storing Logs in MongoDB**
Logs are automatically stored in MongoDB:
- **Normal Logs**: Stored in the `logs` collection.
- **Error Logs**: Stored in the `error_logs` collection.

No additional configuration is needed as this is handled by the `mongoWriter`.

---

### **Important Notes**
1. Ensure the logger is properly initialized in the `main.go` file.
2. Always provide context with logs using the key-value `map[string]interface{}` format for better debugging and traceability.
3. Use the `Info` method for normal events and `Error` for errors.

---

For further assistance, contact the development team lead.

