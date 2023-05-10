SELECT
 tbc.TownshipID
,tbc.[Description]
,ISNULL(SUM(bpic.NewImprAssmt),0) NewImprAssmt
--,SUM(bpic.NewLandAssmt) NewLandAssmt
,ISNULL(SUM(bpic.OldImprAssmt),0) OldImprAssmt
--,SUM(bpic.OldLandAssmt) OldLandAssmt
,ISNULL(SUM(bpic.ImprDiff),0) ImprDiff
--,SUM(bpic.LandDiff) LandDiff
FROM TownshipBoroughCodes tbc
LEFT JOIN byprop_assmtincrease bpic
    ON tbc.TownshipID = bpic.TownShipBorough
WHERE tbc.TownshipID IN
('030',
     '070',
     '150',
     '200',
     '230',
     '020',
     '090',
     '110',
     '130',
     '010',
     '050',
     '170',
     '210',
     '270')
GROUP BY TownshipID, [Description]
ORDER BY TownshipID