package egoscale

// InstancePoolState represents the state of an instance pool.
type InstancePoolState string

const (
	// InstancePoolCreating creating state.
	InstancePoolCreating InstancePoolState = "creating"
	// InstancePoolRunning running state.
	InstancePoolRunning InstancePoolState = "running"
	// InstancePoolDestroying destroying state.
	InstancePoolDestroying InstancePoolState = "destroying"
	// InstancePoolScalingUp scaling up state.
	InstancePoolScalingUp InstancePoolState = "scaling-up"
	// InstancePoolScalingDown scaling down state.
	InstancePoolScalingDown InstancePoolState = "scaling-down"
)

// InstancePool represents an instance pool.
type InstancePool struct {
	ID                *UUID             `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	ServiceOfferingID *UUID             `json:"serviceofferingid"`
	TemplateID        *UUID             `json:"templateid"`
	ZoneID            *UUID             `json:"zoneid"`
	SecurityGroupIDs  []UUID            `json:"securitygroupids"`
	NetworkIDs        []UUID            `json:"networkids"`
	KeyPair           string            `json:"keypair"`
	UserData          string            `json:"userdata"`
	Size              int               `json:"size"`
	RootDiskSize      int               `json:"rootdisksize"`
	State             InstancePoolState `json:"state"`
	VirtualMachines   []VirtualMachine  `json:"virtualmachines"`
}

// CreateInstancePool creates an instance pool.
type CreateInstancePool struct {
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	ServiceOfferingID *UUID  `json:"serviceofferingid"`
	TemplateID        *UUID  `json:"templateid"`
	ZoneID            *UUID  `json:"zoneid"`
	SecurityGroupIDs  []UUID `json:"securitygroupids,omitempty"`
	NetworkIDs        []UUID `json:"networkids,omitempty"`
	KeyPair           string `json:"keypair,omitempty"`
	UserData          string `json:"userdata,omitempty"`
	Size              int    `json:"size"`
	RootDiskSize      int    `json:"rootdisksize,omitempty"`
	_                 bool   `name:"createInstancePool" description:"Create an Instance Pool"`
}

// CreateInstancePoolResponse represents an instance pool creation API response.
type CreateInstancePoolResponse struct {
	ID                *UUID             `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	ServiceOfferingID *UUID             `json:"serviceofferingid"`
	TemplateID        *UUID             `json:"templateid"`
	ZoneID            *UUID             `json:"zoneid"`
	SecurityGroupIDs  []UUID            `json:"securitygroupids"`
	NetworkIDs        []UUID            `json:"networkids"`
	KeyPair           string            `json:"keypair"`
	UserData          string            `json:"userdata"`
	Size              int64             `json:"size"`
	RootDiskSize      int               `json:"rootdisksize"`
	State             InstancePoolState `json:"state"`
}

// Response returns the struct to unmarshal.
func (CreateInstancePool) Response() interface{} {
	return new(CreateInstancePoolResponse)
}

// UpdateInstancePool updates an instance pool.
type UpdateInstancePool struct {
	ID          *UUID  `json:"id"`
	ZoneID      *UUID  `json:"zoneid"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TemplateID  *UUID  `json:"templateid,omitempty"`
	UserData    string `json:"userdata,omitempty"`
	_           bool   `name:"updateInstancePool" description:"Update an Instance Pool"`
}

// UpdateInstancePoolResponse update instance pool response.
type UpdateInstancePoolResponse struct {
	Success bool `json:"success"`
}

func (UpdateInstancePool) Response() interface{} {
	return new(UpdateInstancePoolResponse)
}

// ScaleInstancePool scales an instance pool.
type ScaleInstancePool struct {
	ID     *UUID `json:"id"`
	ZoneID *UUID `json:"zoneid"`
	Size   int   `json:"size"`
	_      bool  `name:"scaleInstancePool" description:"Scale an Instance Pool"`
}

// ScaleInstancePoolResponse scale instance pool response.
type ScaleInstancePoolResponse struct {
	Success bool `json:"success"`
}

func (ScaleInstancePool) Response() interface{} {
	return new(ScaleInstancePoolResponse)
}

// DestroyInstancePool destroys an instance pool.
type DestroyInstancePool struct {
	ID     *UUID `json:"id"`
	ZoneID *UUID `json:"zoneid"`
	_      bool  `name:"destroyInstancePool" description:"Destroy an Instance Pool"`
}

// DestroyInstancePoolResponse destroy instance pool response.
type DestroyInstancePoolResponse struct {
	Success bool `json:"success"`
}

func (DestroyInstancePool) Response() interface{} {
	return new(DestroyInstancePoolResponse)
}

// GetInstancePool retrieves an instance pool's details.
type GetInstancePool struct {
	ID     *UUID `json:"id"`
	ZoneID *UUID `json:"zoneid"`
	_      bool  `name:"getInstancePool" description:"Get an Instance Pool"`
}

// GetInstancePoolResponse get instance pool response.
type GetInstancePoolResponse struct {
	Count         int
	InstancePools []InstancePool `json:"instancepool"`
}

func (GetInstancePool) Response() interface{} {
	return new(GetInstancePoolResponse)
}

// ListInstancePools lists instance pools.
type ListInstancePools struct {
	ZoneID *UUID `json:"zoneid"`
	_      bool  `name:"listInstancePools" description:"List Instance Pools"`
}

// ListInstancePoolsResponse list instance pool response.
type ListInstancePoolsResponse struct {
	Count         int
	InstancePools []InstancePool `json:"instancepool"`
}

func (ListInstancePools) Response() interface{} {
	return new(ListInstancePoolsResponse)
}
