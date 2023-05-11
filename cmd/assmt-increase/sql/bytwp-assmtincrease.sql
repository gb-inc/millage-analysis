SELECT
 tbc.TownshipID
,tbc.[Description]
,ISNULL(SUM(bpic.NewImprAssmt),0) NewImprAssmt
,ISNULL(SUM(bpic.OldImprAssmt),0) OldImprAssmt
,ISNULL(SUM(bpic.ImprDiff),0) ImprDiff
FROM TownshipBoroughCodes tbc
LEFT JOIN byprop_assmtincrease bpic
    ON tbc.TownshipID = bpic.TownShipBorough
GROUP BY TownshipID, [Description]
ORDER BY TownshipID