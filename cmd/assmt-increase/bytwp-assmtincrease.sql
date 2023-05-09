SELECT
 bpic.TownShipBorough
,bpic.DisName
,SUM(bpic.NewImprAssmt) NewImprAssmt
,SUM(bpic.NewLandAssmt) NewLandAssmt
,SUM(bpic.OldImprAssmt) OldImprAssmt
,SUM(bpic.OldLandAssmt) OldLandAssmt
,SUM(bpic.ImprDiff) ImprDiff
,SUM(bpic.LandDiff) LandDiff
FROM byprop_assmtincrease bpic
GROUP BY TownShipBorough, DisName
ORDER BY TownShipBorough