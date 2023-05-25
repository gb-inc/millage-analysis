-- View operating off temporary table storing formal notice values
CREATE OR ALTER VIEW byprop_assmtincrease AS
    SELECT 
    P.TownShipBorough
    ,IIF(TownshipOrBorough = 'B', TBC.[Description] + ' BORO', TBC.[Description] + ' TWP') DisName
    ,P.ControlNumber
    ,AN.ImprovementAssessment AS NewImprAssmt
    ,NoticeVals.IMPR AS OldImprAssmt
    ,(AN.ImprovementAssessment - NoticeVals.IMPR) ImprDiff
    FROM Sheet1$ NoticeVals
    LEFT JOIN Property P
        ON P.ControlNumber = NoticeVals.[CONTROL]
    INNER JOIN AssmtNew AN
        ON AN.PropertyID = P.PropertyID
    INNER JOIN TownshipBoroughCodes TBC
        ON P.TownShipBorough = tbc.TownshipID
    INNER JOIN LandUseCode LUC
        ON P.LandUse = LUC.Code
    WHERE P.ParcelEliminated IS NULL
    AND AN.CurrentAssessment = -1
    AND P.LandUse NOT IN ('166','W10')
    AND TBC.TownshipOrBorough IN ('T','B')
    AND NOT luc.ParcelStatus = 'N'
    AND (AN.ImprovementAssessment) - (NoticeVals.IMPR) >= 0

-- Old View
-- CREATE OR ALTER VIEW byprop_assmtincrease AS
--     SELECT
--      P.TownShipBorough
--     ,IIF(TownshipOrBorough = 'B', TBC.[Description] + ' BORO', TBC.[Description] + ' TWP') DisName
--     ,P.ControlNumber
--     ,AN.ImprovementAssessment AS NewImprAssmt
--     ,AN.LandAssessment AS NewLandAssmt
--     ,ANOld.ImprovementAssessment AS OldImprAssmt
--     ,ANOld.LandAssessment AS OldLandAssmt
--     ,(AN.ImprovementAssessment - ANOLD.ImprovementAssessment) ImprDiff
--     ,(AN.LandAssessment - ANOLD.LandAssessment) LandDiff
--     FROM AssmtNew AN
--     INNER JOIN Property P
--         ON P.PropertyID = AN.PropertyID
--     INNER JOIN TownshipBoroughCodes TBC
--         ON P.TownShipBorough = tbc.TownshipID
--     INNER JOIN LandUseCode LUC
--         ON P.LandUse = LUC.Code
--     LEFT JOIN (
--         SELECT *
--         FROM (
--             SELECT *, RANK() OVER (PARTITION BY PropertyID ORDER BY AssessmentID DESC) AS rank
--             FROM AssmtNew
--             WHERE EffectiveDate >= '2022-07-01 00:00:00.000'
--             AND EffectiveDate <= '2023-01-01 00:00:00.000'
--             AND ApprovalDate >= '2022-02-22 00:00:00.000' --oldest approval date for 07/01 effective dates
-- 		    AND ApprovalDate < '2023-01-01 00:00:00.000'
--         ) ranked
--         WHERE rank = 1
--     ) ANOLD ON P.PropertyID = ANOLD.PropertyID
--     WHERE AN.CurrentAssessment = -1
--     AND AN.ApprovalDate > '2022-07-01 00:00:00.000'
--     AND P.ParcelEliminated IS NULL
--     AND P.LandUse NOT IN ('166','W10')
--     AND TBC.TownshipOrBorough IN ('T','B')
--     AND NOT luc.ParcelStatus = 'N'
--     AND (AN.ImprovementAssessment) - (ANOLD.ImprovementAssessment) >= 0