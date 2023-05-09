SELECT
 bpvc.TownShipBorough
,bpvc.DisName
,SUM(bpvc.NewImprAssmt) NewImprAssmt
,SUM(bpvc.NewLandAssmt) NewLandAssmt
,SUM(bpvc.OldImprAssmt) OldImprAssmt
,SUM(bpvc.OldLandAssmt) OldLandAssmt
,SUM(bpvc.ImprDiff) ImprDiff
,SUM(bpvc.LandDiff) LandDiff
FROM byprop_valuationchange bpvc
GROUP BY TownShipBorough, DisName
ORDER BY TownShipBorough