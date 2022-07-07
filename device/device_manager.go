package device

// 设备管理模块
// 维护一个 sync.Map 结构 key 为设备的唯一sn
//

import "sync"

var ManagerGlabal Manager

type Manager struct {
	devices sync.Map
}

func (m *Manager)Update(device Device) {
	// todo 是否新增，通知上线等
	m.Set(device.Sn, &device)
}

func (m *Manager) List() []Device {
	var devicelist []Device
	m.devices.Range(func(_, dev interface{}) bool {
		device := dev.(*Device)
		devicelist = append(devicelist, *device)
		return true
	})
	return devicelist
}

func (m *Manager) Get(sn string) (interface{}, bool) {
	return m.devices.Load(sn)

}

func (m *Manager) Set(sn string, dev *Device) {

	m.devices.Store(sn, dev)
}

func (m *Manager) Delete(sn string) {
	m.devices.Delete(sn)
}
