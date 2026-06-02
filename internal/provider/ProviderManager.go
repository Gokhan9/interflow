package provider

import "fmt"

// ProviderManager, aktif olan tüm AI providerları yönetmek için kullanılacak.
type ProviderManager struct {
	providers map[string]Provider // Sağlayıcıları adlarına göre saklamak için bir harita
}

func NewManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]Provider), // Sağlayıcı haritası başlatılır
	}
}

// RegisterProvider, yeni bir provider'ı ProviderManager'a eklemek için kullanılır. Provider, "GetName()" metodu ile adını alır ve bu ad, provider'ı saklamak için haritada anahtar olarak kullanılır. Böylece, farklı providerlar kolayca yönetilebilir ve erişilebilir hale gelir.
func (m *ProviderManager) RegisterProvider(p Provider) {
	m.providers[p.GetName()] = p // Sağlayıcı, adını anahtar olarak kullanarak haritaya eklenir
}

func (m *ProviderManager) GetProvider(name string) (Provider, error) {
	p, ok := m.providers[name] // Sağlayıcı adıyla haritadan provider alınır
	if !ok {
		return nil, fmt.Errorf("provider not found %s", name)
	}

	return p, nil // Sağlayıcı başarıyla bulunursa döndürülür
}
