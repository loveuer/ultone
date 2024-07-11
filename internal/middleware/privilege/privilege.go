package privilege

import (
	"fmt"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/resp"
	"github.com/samber/lo"
	"strings"
	"ultone/internal/model"
)

type Relation int64

type vf func(user *model.User, ps ...model.Privilege) error

const (
	RelationAnd Relation = iota + 1
	RelationOr
)

var (
	AndFunc vf = func(user *model.User, ps ...model.Privilege) error {
		pm := lo.SliceToMap(user.Privileges, func(item model.Privilege) (int64, struct{}) {
			return item.Value(), struct{}{}
		})

		for _, p := range ps {
			if _, exist := pm[p.Value()]; !exist {
				return fmt.Errorf("缺少权限: %d, %s, %s", p.Value(), p.Code(), p.Label())
			}
		}

		return nil
	}

	OrFunc vf = func(user *model.User, ps ...model.Privilege) error {
		pm := lo.SliceToMap(user.Privileges, func(item model.Privilege) (int64, struct{}) {
			return item.Value(), struct{}{}
		})

		for _, p := range ps {
			if _, exist := pm[p.Value()]; exist {
				return nil
			}
		}

		return fmt.Errorf("缺少权限: %s", strings.Join(
			lo.Map(ps, func(item model.Privilege, index int) string {
				return item.Code()
			}),
			", ",
		))
	}
)

func Verify(relation Relation, privileges ...model.Privilege) nf.HandlerFunc {

	var _vf vf

	switch relation {
	case RelationAnd:
		_vf = AndFunc
	case RelationOr:
		_vf = OrFunc
	default:
		log.Panic("middleware.Verify: unknown relation")
	}

	return func(c *nf.Ctx) error {
		if len(privileges) == 0 {
			return c.Next()
		}

		op, ok := c.Locals("user").(*model.User)
		if !ok {
			return resp.Resp401(c, nil)
		}

		if err := _vf(op, privileges...); err != nil {
			return resp.Resp403(c, err.Error())
		}

		return c.Next()
	}
}
