# Provider Architecture Diagram

## Current Architecture (Data Sources Only)

```mermaid
graph TB
    subgraph "Terraform Configuration"
        TF[Terraform Config]
    end
    
    subgraph "Provider: registry.terraform.io/morganpeat/environment"
        P[environmentProvider]
        P --> DS1[environment_variable<br/>Data Source]
        P --> DS2[environment_sensitive_variable<br/>Data Source]
    end
    
    subgraph "System"
        ENV[Environment Variables<br/>PATH, HOME, etc.]
    end
    
    TF -->|data block| DS1
    TF -->|data block| DS2
    DS1 -->|os.LookupEnv| ENV
    DS2 -->|os.LookupEnv| ENV
    
    DS1 -->|value| STATE[Terraform State]
    DS2 -->|value sensitive| STATE
    
    style DS1 fill:#e1f5ff
    style DS2 fill:#ffe1e1
    style STATE fill:#fff4e1
```

## Proposed Architecture (Data Sources + Functions)

```mermaid
graph TB
    subgraph "Terraform Configuration"
        TF[Terraform Config]
    end
    
    subgraph "Provider: registry.terraform.io/morganpeat/environment"
        P[environmentProvider]
        
        subgraph "Data Sources Terraform 1.0+"
            DS1[environment_variable<br/>Data Source]
            DS2[environment_sensitive_variable<br/>Data Source]
        end
        
        subgraph "Functions Terraform 1.8+"
            F1[variable<br/>Function]
            F2[sensitive_variable<br/>Function]
        end
        
        P --> DS1
        P --> DS2
        P --> F1
        P --> F2
    end
    
    subgraph "System"
        ENV[Environment Variables<br/>PATH, HOME, API_KEY, etc.]
    end
    
    TF -->|data block| DS1
    TF -->|data block| DS2
    TF -->|provider::environment::variable| F1
    TF -->|provider::environment::sensitive_variable| F2
    
    DS1 -->|os.LookupEnv| ENV
    DS2 -->|os.LookupEnv| ENV
    F1 -->|os.LookupEnv| ENV
    F2 -->|os.LookupEnv| ENV
    
    DS1 -->|cached value| STATE[Terraform State]
    DS2 -->|cached sensitive value| STATE
    F1 -.->|not stored| STATE
    F2 -.->|not stored| STATE
    
    style DS1 fill:#e1f5ff
    style DS2 fill:#ffe1e1
    style F1 fill:#e1ffe1
    style F2 fill:#ffe1f5
    style STATE fill:#fff4e1
```

## Usage Comparison

### Data Source Approach (Terraform 1.0+)

```mermaid
sequenceDiagram
    participant User as Terraform Config
    participant DS as Data Source
    participant OS as Environment
    participant State as Terraform State
    
    User->>DS: data "environment_variable" "path" { name = "PATH" }
    DS->>OS: os.LookupEnv("PATH")
    OS-->>DS: "/usr/bin:/bin"
    DS->>State: Store value
    State-->>User: data.environment_variable.path.value
```

### Function Approach (Terraform 1.8+)

```mermaid
sequenceDiagram
    participant User as Terraform Config
    participant Fn as Function
    participant OS as Environment
    
    User->>Fn: provider::environment::variable("PATH")
    Fn->>OS: os.LookupEnv("PATH")
    OS-->>Fn: "/usr/bin:/bin"
    Fn-->>User: "/usr/bin:/bin"
    Note over Fn,User: No state storage
```

## Implementation Flow

```mermaid
graph LR
    subgraph "Phase 1: Core Implementation"
        A1[Add Functions method<br/>to provider.go]
        A2[Implement<br/>variable_function.go]
        A3[Implement<br/>sensitive_variable_function.go]
        A4[Add unit tests]
        
        A1 --> A2
        A2 --> A3
        A3 --> A4
    end
    
    subgraph "Phase 2: Testing"
        B1[Acceptance tests]
        B2[Integration testing]
        B3[Backward compatibility<br/>verification]
        
        A4 --> B1
        B1 --> B2
        B2 --> B3
    end
    
    subgraph "Phase 3: Documentation"
        C1[Update templates]
        C2[Create examples]
        C3[Generate docs]
        C4[Update README]
        
        B3 --> C1
        C1 --> C2
        C2 --> C3
        C3 --> C4
    end
    
    subgraph "Phase 4: Release"
        D1[Update CHANGELOG]
        D2[Create release notes]
        D3[Tag release]
        
        C4 --> D1
        D1 --> D2
        D2 --> D3
    end
    
    style A1 fill:#e1f5ff
    style A2 fill:#e1f5ff
    style A3 fill:#e1f5ff
    style A4 fill:#e1f5ff
    style B1 fill:#ffe1e1
    style B2 fill:#ffe1e1
    style B3 fill:#ffe1e1
    style C1 fill:#e1ffe1
    style C2 fill:#e1ffe1
    style C3 fill:#e1ffe1
    style C4 fill:#e1ffe1
    style D1 fill:#fff4e1
    style D2 fill:#fff4e1
    style D3 fill:#fff4e1
```

## File Structure

```
terraform-provider-environment/
├── internal/provider/
│   ├── provider.go                          # Updated: Add Functions() method
│   ├── provider_test.go                     # Existing
│   │
│   ├── variable_data_source.go              # Existing: Terraform 1.0+
│   ├── variable_data_source_test.go         # Existing
│   ├── sensitive_variable_data_source.go    # Existing: Terraform 1.0+
│   ├── sensitive_variable_data_source_test.go # Existing
│   │
│   ├── variable_function.go                 # New: Terraform 1.8+
│   ├── variable_function_test.go            # New
│   ├── sensitive_variable_function.go       # New: Terraform 1.8+
│   └── sensitive_variable_function_test.go  # New
│
├── examples/
│   ├── data-sources/                        # Existing examples
│   │   ├── environment_variable/
│   │   └── environment_sensitive_variable/
│   │
│   └── functions/                           # New examples
│       ├── environment_variable/
│       │   └── function.tf
│       └── environment_sensitive_variable/
│           └── function.tf
│
├── docs/
│   ├── data-sources/                        # Existing docs
│   │   ├── variable.md
│   │   └── sensitive_variable.md
│   │
│   └── functions/                           # New docs (auto-generated)
│       ├── variable.md
│       └── sensitive_variable.md
│
├── templates/
│   └── index.md.tmpl                        # Updated: Add functions section
│
├── DESIGN_FUNCTIONS.md                      # This design document
├── ARCHITECTURE_DIAGRAM.md                  # This file
└── README.md                                # Updated: Add functions example
```

## Error Handling Flow

```mermaid
graph TD
    Start[Function Called]
    Start --> Parse[Parse Arguments]
    Parse --> Valid{Arguments Valid?}
    
    Valid -->|No| Error1[Return Argument Error]
    Valid -->|Yes| Lookup[os.LookupEnv]
    
    Lookup --> Found{Variable Found?}
    Found -->|No| Error2[Return Not Found Error]
    Found -->|Yes| Return[Return Value]
    
    Return --> Sensitive{Sensitive Function?}
    Sensitive -->|Yes| MarkSensitive[Mark as Sensitive]
    Sensitive -->|No| Done[Done]
    MarkSensitive --> Done
    
    style Error1 fill:#ffe1e1
    style Error2 fill:#ffe1e1
    style Return fill:#e1ffe1
    style MarkSensitive fill:#ffe1f5
    style Done fill:#e1f5ff
```

## Testing Strategy

```mermaid
graph TB
    subgraph "Unit Tests"
        UT1[Variable exists]
        UT2[Variable not found]
        UT3[Empty value]
        UT4[Special characters]
        UT5[Sensitive marking]
    end
    
    subgraph "Acceptance Tests"
        AT1[Function in output]
        AT2[Function in locals]
        AT3[Function in resource]
        AT4[Sensitive function]
        AT5[Error handling]
    end
    
    subgraph "Integration Tests"
        IT1[Terraform 1.8+]
        IT2[Data sources still work]
        IT3[Mixed usage]
    end
    
    UT1 --> AT1
    UT2 --> AT5
    UT3 --> AT1
    UT4 --> AT1
    UT5 --> AT4
    
    AT1 --> IT1
    AT2 --> IT1
    AT3 --> IT1
    AT4 --> IT1
    AT5 --> IT1
    
    IT1 --> IT2
    IT2 --> IT3
    
    style UT1 fill:#e1f5ff
    style UT2 fill:#e1f5ff
    style UT3 fill:#e1f5ff
    style UT4 fill:#e1f5ff
    style UT5 fill:#e1f5ff
    style AT1 fill:#ffe1e1
    style AT2 fill:#ffe1e1
    style AT3 fill:#ffe1e1
    style AT4 fill:#ffe1e1
    style AT5 fill:#ffe1e1
    style IT1 fill:#e1ffe1
    style IT2 fill:#e1ffe1
    style IT3 fill:#e1ffe1
```

## Backward Compatibility

```mermaid
graph LR
    subgraph "Terraform 1.0 - 1.7"
        T1[User Config]
        T1 --> DS[Data Sources Only]
        DS --> Works1[✓ Works]
    end
    
    subgraph "Terraform 1.8+"
        T2[User Config]
        T2 --> Choice{User Choice}
        Choice -->|Legacy| DS2[Data Sources]
        Choice -->|Modern| FN[Functions]
        Choice -->|Both| BOTH[Mixed Usage]
        
        DS2 --> Works2[✓ Works]
        FN --> Works3[✓ Works]
        BOTH --> Works4[✓ Works]
    end
    
    style Works1 fill:#e1ffe1
    style Works2 fill:#e1ffe1
    style Works3 fill:#e1ffe1
    style Works4 fill:#e1ffe1
```

## Key Design Decisions Summary

| Decision | Rationale |
|----------|-----------|
| Two separate functions (variable, sensitive_variable) | Terraform function return types must be statically defined; cannot conditionally mark as sensitive |
| Keep existing data sources unchanged | Maintain backward compatibility with Terraform < 1.8 |
| Mirror data source naming | Consistency with existing API |
| Error on missing variable | Explicit failure prevents silent configuration errors |
| No state storage for functions | Functions are evaluated on-demand per Terraform design |
| Separate implementation files | Follows existing code organization patterns |
