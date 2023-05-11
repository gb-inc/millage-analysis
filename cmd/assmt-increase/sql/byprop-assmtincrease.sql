CREATE OR ALTER VIEW byprop_assmtincrease AS
    SELECT
     P.TownShipBorough
    ,IIF(TownshipOrBorough = 'B', TBC.[Description] + ' BORO', TBC.[Description] + ' TWP') DisName
    ,P.ControlNumber
    ,AN.ImprovementAssessment AS NewImprAssmt
    ,AN.LandAssessment AS NewLandAssmt
    ,ANOld.ImprovementAssessment AS OldImprAssmt
    ,ANOld.LandAssessment AS OldLandAssmt
    ,(AN.ImprovementAssessment - ANOLD.ImprovementAssessment) ImprDiff
    ,(AN.LandAssessment - ANOLD.LandAssessment) LandDiff
    FROM AssmtNew AN
    INNER JOIN Property P
        ON P.PropertyID = AN.PropertyID
    INNER JOIN TownshipBoroughCodes TBC
        ON P.TownShipBorough = tbc.TownshipID
    INNER JOIN LandUseCode LUC
        ON P.LandUse = LUC.Code
    LEFT JOIN (
        SELECT *
        FROM (
            SELECT *, RANK() OVER (PARTITION BY PropertyID ORDER BY AssessmentID DESC) AS rank
            FROM AssmtNew
            WHERE EffectiveDate = '2023-01-01 00:00:00.000'
            AND ApprovalDate < '2023-01-01 00:00:00.000'
        ) ranked
        WHERE rank = 1
    ) ANOLD ON P.PropertyID = ANOLD.PropertyID
    WHERE AN.CurrentAssessment = -1
    AND AN.ApprovalDate > '2023-01-01 00:00:00.000'
    AND P.ParcelEliminated IS NULL
    AND P.LandUse NOT IN ('166','W10')
    AND TBC.TownshipOrBorough IN ('T','B')
    AND NOT luc.ParcelStatus = 'N'
    AND (AN.ImprovementAssessment) - (ANOLD.ImprovementAssessment) >= 0
    AND P.TownShipBorough IN