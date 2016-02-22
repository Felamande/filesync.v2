package pairs

import(
    // "github.com/lunny/tango"
    // "errors"
    "github.com/Felamande/filesync.v2/syncer"
    "github.com/Felamande/filesync.v2/server/routers/base"
)

type GetAllRouter struct{
    base.BaseJSONRouter
    
}

func (r *GetAllRouter)Get()interface{}{
    return syncer.Default().PairMap
}
