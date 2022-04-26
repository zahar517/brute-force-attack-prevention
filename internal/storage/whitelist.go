package storage

import (
	"context"
)

func (s *Storage) CreateWhitelistSubnet(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO whitelist (subnet) VALUES ($1)`, subnet)

	return err
}

func (s *Storage) DeleteWhitelistSubnet(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM whitelist WHERE subnet = $1`, subnet)

	return err
}

func (s *Storage) IsIpInWhitelist(ctx context.Context, ip string) (bool, error) {
	res, err := s.db.ExecContext(ctx, `SELECT subnet FROM whitelist WHERE subnet >>= $1`, ip)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if n == 0 {
		return false, nil
	}

	return true, nil
}
